//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package identity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/ericls/imgdd/db/.gen/imgdd/public/model"
	. "github.com/ericls/imgdd/db/.gen/imgdd/public/table"
	dm "github.com/ericls/imgdd/domainmodels"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

// APIKeySecretPrefix marks every API key secret so leaked keys are easy to
// recognize (e.g. by secret scanners).
const APIKeySecretPrefix = "imgdd_"

// apiKeyDisplayPrefixLen is how much of the secret is stored in plain text so
// users can tell their keys apart.
const apiKeyDisplayPrefixLen = len(APIKeySecretPrefix) + 6

const APIKeyMaxNameLength = 100

var ErrAPIKeyNameInvalid = errors.New("api key name must be between 1 and 100 characters")
var ErrAPIKeyNotFound = errors.New("api key not found")

func generateAPIKeySecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return APIKeySecretPrefix + base64.RawURLEncoding.EncodeToString(b), nil
}

// HashAPIKeySecret returns the value stored in api_key_table.key_hash.
// Secrets carry 256 bits of randomness, so a fast hash is sufficient.
func HashAPIKeySecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

func convertAPIKey(k *model.APIKeyTable) *dm.APIKey {
	return &dm.APIKey{
		Id:                 k.ID.String(),
		OrganizationUserId: k.OrganizationUserID.String(),
		Name:               k.Name,
		KeyPrefix:          k.KeyPrefix,
		LastUsedAt:         k.LastUsedAt,
		RevokedAt:          k.RevokedAt,
		CreatedAt:          k.CreatedAt,
	}
}

// CreateAPIKey creates a key for the organization user and returns it along
// with the secret. The secret is only available here; it cannot be recovered.
func (repo *DBIdentityRepo) CreateAPIKey(organizationUserId string, name string) (*dm.APIKey, string, error) {
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > APIKeyMaxNameLength {
		return nil, "", ErrAPIKeyNameInvalid
	}
	orgUserUUID, err := uuid.Parse(organizationUserId)
	if err != nil {
		return nil, "", err
	}
	secret, err := generateAPIKeySecret()
	if err != nil {
		return nil, "", err
	}
	stmt := APIKeyTable.INSERT(
		APIKeyTable.OrganizationUserID,
		APIKeyTable.Name,
		APIKeyTable.KeyPrefix,
		APIKeyTable.KeyHash,
	).VALUES(
		orgUserUUID, name, secret[:apiKeyDisplayPrefixLen], HashAPIKeySecret(secret),
	).RETURNING(APIKeyTable.AllColumns)
	dest := model.APIKeyTable{}
	if err := stmt.Query(repo.DB, &dest); err != nil {
		return nil, "", err
	}
	return convertAPIKey(&dest), secret, nil
}

// ListAPIKeys returns the organization user's keys that have not been revoked,
// newest first.
func (repo *DBIdentityRepo) ListAPIKeys(organizationUserId string) ([]*dm.APIKey, error) {
	orgUserUUID, err := uuid.Parse(organizationUserId)
	if err != nil {
		return nil, err
	}
	stmt := SELECT(APIKeyTable.AllColumns).FROM(APIKeyTable).WHERE(
		APIKeyTable.OrganizationUserID.EQ(UUID(orgUserUUID)).
			AND(APIKeyTable.RevokedAt.IS_NULL()),
	).ORDER_BY(APIKeyTable.CreatedAt.DESC())
	dest := []model.APIKeyTable{}
	if err := stmt.Query(repo.DB, &dest); err != nil {
		return nil, err
	}
	ret := make([]*dm.APIKey, len(dest))
	for i := range dest {
		ret[i] = convertAPIKey(&dest[i])
	}
	return ret, nil
}

// RevokeAPIKey revokes a key owned by the organization user. Revoking a key
// that does not exist, belongs to someone else, or is already revoked returns
// ErrAPIKeyNotFound.
func (repo *DBIdentityRepo) RevokeAPIKey(organizationUserId string, apiKeyId string) error {
	orgUserUUID, err := uuid.Parse(organizationUserId)
	if err != nil {
		return ErrAPIKeyNotFound
	}
	keyUUID, err := uuid.Parse(apiKeyId)
	if err != nil {
		return ErrAPIKeyNotFound
	}
	now := time.Now()
	stmt := APIKeyTable.UPDATE(APIKeyTable.RevokedAt, APIKeyTable.UpdatedAt).
		SET(TimestampzT(now), TimestampzT(now)).
		WHERE(
			APIKeyTable.ID.EQ(UUID(keyUUID)).
				AND(APIKeyTable.OrganizationUserID.EQ(UUID(orgUserUUID))).
				AND(APIKeyTable.RevokedAt.IS_NULL()),
		)
	res, err := stmt.Exec(repo.DB)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrAPIKeyNotFound
	}
	return nil
}

// GetOrganizationUserByAPIKeySecret resolves a secret to the organization user
// that owns it, or nil if the secret is unknown or revoked. It also records
// when the key was last used.
func (repo *DBIdentityRepo) GetOrganizationUserByAPIKeySecret(secret string) *dm.OrganizationUser {
	if !strings.HasPrefix(secret, APIKeySecretPrefix) {
		return nil
	}
	stmt := APIKeyTable.UPDATE(APIKeyTable.LastUsedAt).
		SET(TimestampzT(time.Now())).
		WHERE(
			APIKeyTable.KeyHash.EQ(String(HashAPIKeySecret(secret))).
				AND(APIKeyTable.RevokedAt.IS_NULL()),
		).RETURNING(APIKeyTable.OrganizationUserID)
	dest := model.APIKeyTable{}
	if err := stmt.Query(repo.DB, &dest); err != nil {
		return nil
	}
	return repo.GetOrganizationUserById(dest.OrganizationUserID.String())
}
