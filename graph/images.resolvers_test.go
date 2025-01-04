package graph_test

import (
	"imgdd/domainmodels"
	"imgdd/graph/model"
	"imgdd/utils"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func mockSaveFunc(file utils.SeekerReader, filename string, mimeType string) error {
	return nil
}

func createStorageDefinition(t *testing.T, tc *TestContext) *domainmodels.StorageDefinition {
	configJSON := `{
		"bucket": "test-bucket",
		"endpoint": "us-west-2",
		"access": "test",
		"secret": "test"
	}`
	tc.storageRepo.CreateStorageDefinition("s3", configJSON, "test1", true, 1)
	sd, err := tc.storageRepo.GetStorageDefinitionByIdentifier("test1")
	require.NoError(t, err)
	return sd
}

func createImage(t *testing.T, tc *TestContext, uploaderId string, storageDefId string) *domainmodels.Image {
	identifier := uuid.New().String()
	fakeImage := domainmodels.Image{
		UploaderIP:      "127.0.0.1",
		CreatedById:     uploaderId,
		MIMEType:        "image/png",
		Name:            identifier + ".png",
		Identifier:      identifier,
		NominalByteSize: int32(100),
		NominalWidth:    100,
		NominalHeight:   100,
	}
	si, err := tc.imageRepo.CreateAndSaveUploadedImage(&fakeImage, []byte(""), storageDefId, mockSaveFunc)
	require.NoError(t, err)
	image := si.Image
	require.Equal(t, fakeImage.Name, image.Name)
	return image
}

func tImagesNoFilterNoOrderSiteOwner(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			Images *model.ImagesResult
		}
	}
	orgUser := tc.forceAuthenticate(asSiteOwner)
	sd := createStorageDefinition(t, tc)
	createImage(t, tc, orgUser.Id, sd.Id)
	err := tc.client.Post(`
	query foo {
		viewer {
			images {
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
					totalCount
				}
				edges {
					cursor
					node {
						id
						url
					}
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.NotNil(t, resp.Viewer.Images)
	require.Len(t, resp.Viewer.Images.Edges, 1)
	require.NotNil(t, resp.Viewer.Images.PageInfo)
	require.False(t, resp.Viewer.Images.PageInfo.HasNextPage)
	require.False(t, resp.Viewer.Images.PageInfo.HasPreviousPage)
	require.NotNil(t, resp.Viewer.Images.PageInfo.StartCursor)
	imageEdge := resp.Viewer.Images.Edges[0]
	require.NotNil(t, imageEdge)
	require.NotEmpty(t, imageEdge.Node.ID)
	require.NotEmpty(t, imageEdge.Node.URL)
	require.Equal(t, resp.Viewer.Images.PageInfo.StartCursor, &imageEdge.Cursor)
}

func tSiteOwnerCanAccessAllImages(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			Images *model.ImagesResult
		}
	}
	orgUser := tc.forceAuthenticate(asSiteOwner)
	sd := createStorageDefinition(t, tc)
	createImage(t, tc, orgUser.Id, sd.Id)
	createImage(t, tc, "", sd.Id)
	err := tc.client.Post(`
	query foo {
		viewer {
			images {
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
					totalCount
				}
				edges {
					cursor
					node {
						id
						url
					}
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.NotNil(t, resp.Viewer.Images)
	require.Len(t, resp.Viewer.Images.Edges, 2)
}

func tNormalUserCanOnlyAcessOwnImages(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			Images *model.ImagesResult
		}
	}
	orgUser := tc.forceAuthenticate()
	sd := createStorageDefinition(t, tc)
	createImage(t, tc, orgUser.Id, sd.Id)
	createImage(t, tc, "", sd.Id)
	err := tc.client.Post(`
	query foo {
		viewer {
			images {
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
					totalCount
				}
				edges {
					cursor
					node {
						id
						url
					}
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.NotNil(t, resp.Viewer.Images)
	require.Len(t, resp.Viewer.Images.Edges, 1)
}

func tBasicPagination(t *testing.T, tc *TestContext) {
	type Resp struct {
		Viewer *struct {
			Images *model.ImagesResult
		}
	}
	var resp Resp
	orgUser := tc.forceAuthenticate(asSiteOwner)
	sd := createStorageDefinition(t, tc)
	for i := 0; i < domainmodels.ImageResultPerPage+6; i++ {
		createImage(t, tc, orgUser.Id, sd.Id)
	}
	err := tc.client.Post(`
	query foo {
		viewer {
			images(orderBy: {name: desc}) {
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
					currentCount
				}
				edges {
					cursor
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.NotNil(t, resp.Viewer.Images)
	require.Len(t, resp.Viewer.Images.Edges, domainmodels.ImageResultPerPage)
	require.Equal(t, *resp.Viewer.Images.PageInfo.CurrentCount, domainmodels.ImageResultPerPage)
	require.True(t, resp.Viewer.Images.PageInfo.HasNextPage)
	require.NotNil(t, resp.Viewer.Images.PageInfo.EndCursor)
	require.False(t, resp.Viewer.Images.PageInfo.HasPreviousPage)
	endCursor := *resp.Viewer.Images.PageInfo.EndCursor
	// Get the next page
	var resp2 Resp
	err = tc.client.Post(`
	query foo($after: String) {
		viewer {
			images(orderBy: {name: desc} after: $after) {
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
					totalCount
					currentCount
				}
				edges {
					cursor
					node {
						id
						url
						name
					}
				}
			}
		}
	}`, &resp2, client.Var("after", endCursor))
	require.NoError(t, err)
	require.NotNil(t, resp2.Viewer)
	require.NotNil(t, resp2.Viewer.Images)
	require.Len(t, resp2.Viewer.Images.Edges, 6)
	require.False(t, resp2.Viewer.Images.PageInfo.HasNextPage)
	require.NotNil(t, resp2.Viewer.Images.PageInfo.EndCursor)
	require.True(t, resp2.Viewer.Images.PageInfo.HasPreviousPage)
}

func tBasicPaginationByCreatedAt(t *testing.T, tc *TestContext) {
	type Resp struct {
		Viewer *struct {
			Images *model.ImagesResult
		}
	}
	var resp Resp
	orgUser := tc.forceAuthenticate(asSiteOwner)
	sd := createStorageDefinition(t, tc)
	for i := 0; i < domainmodels.ImageResultPerPage+6; i++ {
		createImage(t, tc, orgUser.Id, sd.Id)
	}
	err := tc.client.Post(`
	query foo {
		viewer {
			images {
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
					currentCount
				}
				edges {
					cursor
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.NotNil(t, resp.Viewer.Images)
	require.Len(t, resp.Viewer.Images.Edges, domainmodels.ImageResultPerPage)
	require.Equal(t, *resp.Viewer.Images.PageInfo.CurrentCount, domainmodels.ImageResultPerPage)
	require.True(t, resp.Viewer.Images.PageInfo.HasNextPage)
	require.NotNil(t, resp.Viewer.Images.PageInfo.EndCursor)
	require.False(t, resp.Viewer.Images.PageInfo.HasPreviousPage)
	endCursor := *resp.Viewer.Images.PageInfo.EndCursor
	// Get the next page
	var resp2 Resp
	err = tc.client.Post(`
	query foo($after: String) {
		viewer {
			images(orderBy: {createdAt: desc} after: $after) {
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
					totalCount
					currentCount
				}
				edges {
					cursor
					node {
						id
						url
						name
					}
				}
			}
		}
	}`, &resp2, client.Var("after", endCursor))
	require.NoError(t, err)
	require.NotNil(t, resp2.Viewer)
	require.NotNil(t, resp2.Viewer.Images)
	require.Len(t, resp2.Viewer.Images.Edges, 6)
	require.False(t, resp2.Viewer.Images.PageInfo.HasNextPage)
	require.NotNil(t, resp2.Viewer.Images.PageInfo.EndCursor)
	require.True(t, resp2.Viewer.Images.PageInfo.HasPreviousPage)
}

func TestImageResolvers(t *testing.T) {
	tc := newTestContext(t)
	tc.runTestCases(
		tImagesNoFilterNoOrderSiteOwner,
		tSiteOwnerCanAccessAllImages,
		tNormalUserCanOnlyAcessOwnImages,
		tBasicPagination,
		tBasicPaginationByCreatedAt,
	)
}
