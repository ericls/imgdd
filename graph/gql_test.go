package graph_test

import (
	"testing"

	"github.com/ericls/imgdd/graph/model"
	"github.com/ericls/imgdd/utils"

	"github.com/stretchr/testify/require"
)

func tAuthenticate(t *testing.T, tc *TestContext) {
	var resp struct {
		Authenticate *model.ViewerResult
	}
	orgUser, err := tc.identityRepo.CreateUserWithOrganization("test@example.com", "test_org", "password")
	if err != nil {
		t.Fatal(err)
	}
	err = tc.client.Post(`
	mutation {
		authenticate(email: "test@example.com", password: "password") {
			viewer {
				id
				organizationUser {
					id
					user {
						id
					}
					organization {
						id
					}
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Authenticate)
	require.Equal(t, orgUser.Id, resp.Authenticate.Viewer.OrganizationUser.ID)
}

func tCreateStorageDefinition(t *testing.T, tc *TestContext) {
	var resp struct {
		CreateStorageDefinition *struct {
			Id         string
			Identifier string
			Config     model.S3StorageConfig
			IsEnabled  bool
			Priority   int
		}
	}

	configJSON := `{
		"bucket": "test-bucket",
		"endpoint": "us-west-2",
		"access": "test",
		"secret": "test"
	}`
	tc.forceAuthenticate(asSiteOwner)
	err := tc.client.Post(`
	mutation {
		createStorageDefinition(input: {
			storageType: S3
			configJSON: "`+utils.JsonEscape(configJSON)+`"
			identifier: "test"
			isEnabled: true
			priority: 1
		}) {
			id
			identifier
			isEnabled
			priority
			config {
				... on S3StorageConfig {
					bucket
					endpoint
					access
					secret
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.CreateStorageDefinition)
	// check that the storage definition was created
	storageDefinition, err := tc.storageDefRepo.GetStorageDefinitionByIdentifier("test")
	require.NoError(t, err)
	require.NotNil(t, storageDefinition)
}

func tCreateStorageDefinitionWithInvalidConfig(t *testing.T, tc *TestContext) {
	var resp struct {
		CreateStorageDefinition *model.StorageDefinition
	}
	tc.forceAuthenticate(asSiteOwner)
	configJSON := `{
		"bucket": "test-bucket",
		"endpoint": "us-west-2",
		"INVALID": "test",
		"secret": "test"
	}`
	err := tc.client.Post(`
	mutation {
		createStorageDefinition(input: {
			storageType: S3
			configJSON: "`+utils.JsonEscape(configJSON)+`"
			identifier: "test"
			isEnabled: true
			priority: 1
		}) {
			id
			identifier
			isEnabled
			priority
			config {
				... on S3StorageConfig {
					bucket
					endpoint
					access
					secret
				}
			}
		}
	}`, &resp)
	require.Error(t, err)
	require.Nil(t, resp.CreateStorageDefinition)
}

func tListStorageDefinitions(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			Id                 string
			StorageDefinitions []*struct {
				Id         string
				Identifier string
				IsEnabled  bool
				Priority   int
				Config     *model.S3StorageConfig
			}
		}
	}
	tc.forceAuthenticate(asSiteOwner)
	gqlQuery := `
	query {
		viewer {
			id
			storageDefinitions {
				id
				identifier
				isEnabled
				priority
				config {
					... on S3StorageConfig {
						bucket
						endpoint
						access
						secret
					}
				}
			}
		}
	}`
	err := tc.client.Post(gqlQuery, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.Len(t, resp.Viewer.StorageDefinitions, 0)
	// create a storage definition
	configJSON := `{
		"bucket": "test-bucket",
		"endpoint": "us-west-2",
		"access": "test",
		"secret": "test"
	}`
	tc.storageDefRepo.CreateStorageDefinition("s3", configJSON, "test1", true, 1)
	err = tc.client.Post(gqlQuery, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.Len(t, resp.Viewer.StorageDefinitions, 1)
}

func tUpdateStorageDefinition(t *testing.T, tc *TestContext) {
	var resp struct {
		UpdateStorageDefinition *model.StorageDefinition
	}
	tc.forceAuthenticate(asSiteOwner)
	configJSON := `{
		"bucket": "test-bucket",
		"endpoint": "us-west-2",
		"access": "test",
		"secret": "test"
	}`
	tc.storageDefRepo.CreateStorageDefinition("s3", configJSON, "test1", true, 1)
	err := tc.client.Post(`
	mutation {
		updateStorageDefinition(input: {
			identifier: "test1"
			isEnabled: false
			priority: 2
		}) {
			id
			identifier
			isEnabled
			priority
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.UpdateStorageDefinition)
	require.Equal(t, false, resp.UpdateStorageDefinition.IsEnabled)
	require.Equal(t, 2, resp.UpdateStorageDefinition.Priority)
}

func TestResolver(t *testing.T) {
	tc := newTestContext(t)
	tc.runTestCases(
		tAuthenticate,
		tCreateStorageDefinition,
		tCreateStorageDefinitionWithInvalidConfig,
		tListStorageDefinitions,
		tUpdateStorageDefinition,
	)
}
