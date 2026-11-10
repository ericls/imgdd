package identity

import (
	"strings"
	"testing"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/test_support"
	"github.com/stretchr/testify/require"
)

func TestAPIKeyLifecycle(t *testing.T) {
	dbConfig := TestServiceMan.GetDBConfig()
	test_support.ResetDatabase(dbConfig)
	repo := NewDBIdentityRepo(db.GetConnection(dbConfig))

	owner, err := repo.CreateUserWithOrganization("keys@home.arpa", "keys", "pw")
	require.NoError(t, err)
	other, err := repo.CreateUserWithOrganization("other@home.arpa", "other", "pw")
	require.NoError(t, err)

	_, _, err = repo.CreateAPIKey(owner.Id, "   ")
	require.ErrorIs(t, err, ErrAPIKeyNameInvalid)
	_, _, err = repo.CreateAPIKey(owner.Id, strings.Repeat("x", APIKeyMaxNameLength+1))
	require.ErrorIs(t, err, ErrAPIKeyNameInvalid)

	key, secret, err := repo.CreateAPIKey(owner.Id, " ci ")
	require.NoError(t, err)
	require.Equal(t, "ci", key.Name)
	require.True(t, strings.HasPrefix(secret, APIKeySecretPrefix))
	require.True(t, strings.HasPrefix(secret, key.KeyPrefix))
	require.Less(t, len(key.KeyPrefix), len(secret))
	require.Nil(t, key.LastUsedAt)

	resolved := repo.GetOrganizationUserByAPIKeySecret(secret)
	require.NotNil(t, resolved)
	require.Equal(t, owner.Id, resolved.Id)
	require.Equal(t, owner.User.Email, resolved.User.Email)

	keys, err := repo.ListAPIKeys(owner.Id)
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.NotNil(t, keys[0].LastUsedAt)

	require.Nil(t, repo.GetOrganizationUserByAPIKeySecret(secret+"x"))
	require.Nil(t, repo.GetOrganizationUserByAPIKeySecret("not-a-key"))

	// Another user can neither see nor revoke the key.
	otherKeys, err := repo.ListAPIKeys(other.Id)
	require.NoError(t, err)
	require.Empty(t, otherKeys)
	require.ErrorIs(t, repo.RevokeAPIKey(other.Id, key.Id), ErrAPIKeyNotFound)
	require.NotNil(t, repo.GetOrganizationUserByAPIKeySecret(secret))

	require.NoError(t, repo.RevokeAPIKey(owner.Id, key.Id))
	require.Nil(t, repo.GetOrganizationUserByAPIKeySecret(secret))
	require.ErrorIs(t, repo.RevokeAPIKey(owner.Id, key.Id), ErrAPIKeyNotFound)
	require.ErrorIs(t, repo.RevokeAPIKey(owner.Id, "not-a-uuid"), ErrAPIKeyNotFound)
	keys, err = repo.ListAPIKeys(owner.Id)
	require.NoError(t, err)
	require.Empty(t, keys)
}
