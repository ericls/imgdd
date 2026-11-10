package graph_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/require"
)

type apiKeyResp struct {
	ID         string
	Name       string
	KeyPrefix  string
	LastUsedAt *string
}

func listAPIKeys(t *testing.T, tc *TestContext) []apiKeyResp {
	t.Helper()
	var resp struct {
		Viewer struct {
			APIKeys []apiKeyResp
		}
	}
	err := tc.client.Post(`query { viewer { apiKeys { id name keyPrefix lastUsedAt } } }`, &resp)
	require.NoError(t, err)
	return resp.Viewer.APIKeys
}

func createAPIKey(tc *TestContext, name string) (apiKeyResp, string, error) {
	var resp struct {
		CreateAPIKey struct {
			APIKey apiKeyResp
			Secret string
		}
	}
	err := tc.client.Post(`
	mutation($name: String!) {
		createAPIKey(input: {name: $name}) {
			apiKey { id name keyPrefix lastUsedAt }
			secret
		}
	}`, &resp, client.Var("name", name))
	return resp.CreateAPIKey.APIKey, resp.CreateAPIKey.Secret, err
}

func revokeAPIKey(tc *TestContext, id string) error {
	var resp struct {
		RevokeAPIKey struct{ ID string }
	}
	return tc.client.Post(`
	mutation($id: ID!) {
		revokeAPIKey(input: {id: $id}) { id }
	}`, &resp, client.Var("id", id))
}

func tAPIKeyLifecycle(t *testing.T, tc *TestContext) {
	owner := tc.forceAuthenticate()
	require.Empty(t, listAPIKeys(t, tc))

	key, secret, err := createAPIKey(tc, "ci")
	require.NoError(t, err)
	require.Equal(t, "ci", key.Name)
	require.NotEmpty(t, secret)
	require.Contains(t, secret, key.KeyPrefix)

	keys := listAPIKeys(t, tc)
	require.Len(t, keys, 1)
	require.Equal(t, key.ID, keys[0].ID)

	resolved := tc.identityRepo.GetOrganizationUserByAPIKeySecret(secret)
	require.NotNil(t, resolved)
	require.Equal(t, owner.Id, resolved.Id)

	require.NoError(t, revokeAPIKey(tc, key.ID))
	require.Empty(t, listAPIKeys(t, tc))
	require.Nil(t, tc.identityRepo.GetOrganizationUserByAPIKeySecret(secret))
	require.Error(t, revokeAPIKey(tc, key.ID))
}

func tAPIKeyInvalidName(t *testing.T, tc *TestContext) {
	tc.forceAuthenticate()
	_, _, err := createAPIKey(tc, "  ")
	require.Error(t, err)
}

func tAPIKeyCannotRevokeOthersKey(t *testing.T, tc *TestContext) {
	tc.forceAuthenticate()
	key, secret, err := createAPIKey(tc, "mine")
	require.NoError(t, err)

	tc.forceAuthenticate()
	require.Empty(t, listAPIKeys(t, tc))
	require.Error(t, revokeAPIKey(tc, key.ID))
	require.NotNil(t, tc.identityRepo.GetOrganizationUserByAPIKeySecret(secret))
}

func tAPIKeyRequiresAuthentication(t *testing.T, tc *TestContext) {
	_, _, err := createAPIKey(tc, "anon")
	require.Error(t, err)
	var resp struct{}
	err = tc.client.Post(`query { viewer { apiKeys { id } } }`, &resp)
	require.Error(t, err)
}

func TestAPIKeyResolvers(t *testing.T) {
	tc := newTestContext(t)
	tc.runTestCases(
		tAPIKeyLifecycle,
		tAPIKeyInvalidName,
		tAPIKeyCannotRevokeOthersKey,
		tAPIKeyRequiresAuthentication,
	)
}
