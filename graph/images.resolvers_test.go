package graph_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/graph/model"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/utils"

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
	tc.storageDefRepo.CreateStorageDefinition("s3", configJSON, "test1", true, 1)
	sd, err := tc.storageDefRepo.GetStorageDefinitionByIdentifier("test1")
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
	si, err := tc.imageRepo.CreateAndSaveUploadedImage(&fakeImage, "image/png", []byte(""), storageDefId, mockSaveFunc)
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

func tDeletingImage(t *testing.T, tc *TestContext) {
	graphqlDoc := `
  mutation foo($input: DeleteImageInput!) {
    deleteImage(input: $input) {
      id
    }
  }`

	orgUser1 := tc.forceAuthenticate()
	orgUser2 := tc.forceAuthenticate()
	sd := createStorageDefinition(t, tc)
	orgUser1Image := createImage(t, tc, orgUser1.Id, sd.Id)
	orgUser2Image := createImage(t, tc, orgUser2.Id, sd.Id)

	// Test that normal user can not delete other user's images
	var resp1 struct {
		DeleteImage *model.DeleteImageResult
	}
	tc.setAuthenticatedUser(orgUser1)
	err := tc.client.Post(graphqlDoc, &resp1, client.Var("input", model.DeleteImageInput{ID: orgUser2Image.Id}))
	require.Error(t, err)
	require.Nil(t, resp1.DeleteImage)
	var resp2 struct {
		DeleteImage *model.DeleteImageResult
	}

	// Test that normal user can delete their own images
	tc.setAuthenticatedUser(orgUser1)
	err = tc.client.Post(graphqlDoc, &resp2, client.Var("input", model.DeleteImageInput{ID: orgUser1Image.Id}))
	require.NoError(t, err)
	require.NotNil(t, resp2.DeleteImage)
	require.Equal(t, orgUser1Image.Id, *resp2.DeleteImage.ID)

	// Test that site owner can delete any image
	orgUser3 := tc.forceAuthenticate(asSiteOwner)
	tc.setAuthenticatedUser(orgUser3)
	var resp3 struct {
		DeleteImage *model.DeleteImageResult
	}
	err = tc.client.Post(graphqlDoc, &resp3, client.Var("input", model.DeleteImageInput{ID: orgUser2Image.Id}))
	require.NoError(t, err)
	require.NotNil(t, resp3.DeleteImage)
	require.Equal(t, orgUser2Image.Id, *resp3.DeleteImage.ID)

	// Test images cannot be retrieved after deletion
	i, _ := tc.imageRepo.GetImageById(orgUser1Image.Id)
	require.Nil(t, i)
	i, _ = tc.imageRepo.GetImageById(orgUser2Image.Id)
	require.Nil(t, i)

	// Test unauthorized user can not delete any image
	image3 := createImage(t, tc, "", sd.Id)
	tc.clearAuthenticationInfo()
	var resp4 struct {
		DeleteImage *model.DeleteImageResult
	}
	err = tc.client.Post(graphqlDoc, &resp4, client.Var("input", model.DeleteImageInput{ID: image3.Id}))
	require.Error(t, err)
}

func tImageCreatedByIsPopulated(t *testing.T, tc *TestContext) {
	orgUser := tc.forceAuthenticate(asSiteOwner)
	sd := createStorageDefinition(t, tc)
	createImage(t, tc, orgUser.Id, sd.Id)

	var resp struct {
		Viewer *struct {
			Images *struct {
				Edges []struct {
					Node struct {
						ID        string
						CreatedBy *struct {
							ID   string
							User struct {
								ID    string
								Email string
							}
						}
					}
				}
			}
		}
	}
	err := tc.client.Post(`
	query {
		viewer {
			images {
				edges {
					node {
						id
						createdBy {
							id
							user {
								id
								email
							}
						}
					}
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.Len(t, resp.Viewer.Images.Edges, 1)
	node := resp.Viewer.Images.Edges[0].Node
	require.NotNil(t, node.CreatedBy, "createdBy must not be nil")
	require.Equal(t, orgUser.Id, node.CreatedBy.ID)
	require.Equal(t, orgUser.User.Email, node.CreatedBy.User.Email)
}

func tImageCreatedByNullWhenNoCreator(t *testing.T, tc *TestContext) {
	tc.forceAuthenticate(asSiteOwner)
	sd := createStorageDefinition(t, tc)
	// empty string createdById → NULL in DB
	createImage(t, tc, "", sd.Id)

	var resp struct {
		Viewer *struct {
			Images *struct {
				Edges []struct {
					Node struct {
						ID        string
						CreatedBy *struct{ ID string }
					}
				}
			}
		}
	}
	err := tc.client.Post(`
	query {
		viewer {
			images {
				edges {
					node {
						id
						createdBy {
							id
						}
					}
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.Len(t, resp.Viewer.Images.Edges, 1)
	require.Nil(t, resp.Viewer.Images.Edges[0].Node.CreatedBy)
}

func makeTestPNGBytes(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func createFSStorageDefinition(t *testing.T, tc *TestContext) (*domainmodels.StorageDefinition, string) {
	tempDir, err := os.MkdirTemp("", "imgdd_test_*")
	require.NoError(t, err)
	configJSON := `{"mediaRoot": "` + tempDir + `"}`
	sd, err := tc.storageDefRepo.CreateStorageDefinition("fs", configJSON, "test-fs", true, 1)
	require.NoError(t, err)
	return sd, tempDir
}

func createRealImage(t *testing.T, tc *TestContext, uploaderId string, storageDefId string, imgBytes []byte) *domainmodels.Image {
	identifier := uuid.New().String()
	fakeImage := domainmodels.Image{
		UploaderIP:      "127.0.0.1",
		CreatedById:     uploaderId,
		MIMEType:        "image/png",
		Name:            identifier + ".png",
		Identifier:      identifier,
		NominalByteSize: int32(len(imgBytes)),
		NominalWidth:    100,
		NominalHeight:   100,
	}
	storageInstance, err := storage.GetStorage(&domainmodels.StorageDefinition{
		Id:          storageDefId,
		StorageType: "fs",
		Config:      func() string { sd, _ := tc.storageDefRepo.GetStorageDefinitionById(storageDefId); return sd.Config }(),
		IsEnabled:   true,
	})
	require.NoError(t, err)
	si, err := tc.imageRepo.CreateAndSaveUploadedImage(&fakeImage, "image/png", imgBytes, storageDefId, storageInstance.Save)
	require.NoError(t, err)
	return si.Image
}

func tApplyWatermark(t *testing.T, tc *TestContext) {
	orgUser := tc.forceAuthenticate()
	sd, tempDir := createFSStorageDefinition(t, tc)
	defer os.RemoveAll(tempDir)

	baseBytes := makeTestPNGBytes(200, 200, color.RGBA{255, 0, 0, 255})
	overlayBytes := makeTestPNGBytes(50, 50, color.RGBA{0, 0, 255, 255})
	baseImage := createRealImage(t, tc, orgUser.Id, sd.Id, baseBytes)
	overlayImage := createRealImage(t, tc, orgUser.Id, sd.Id, overlayBytes)

	var resp struct {
		ApplyWatermark *struct {
			Image *struct {
				ID       string
				Name     string
				MIMEType string
				Parent   *struct {
					ID string
				}
				Changes *string
			}
		}
	}

	err := tc.client.Post(`
	mutation applyWatermark($input: ApplyWatermarkInput!) {
		applyWatermark(input: $input) {
			image {
				id
				name
				MIMEType
				parent {
					id
				}
				changes
			}
		}
	}`, &resp, client.Var("input", map[string]interface{}{
		"baseImageId":    baseImage.Id,
		"overlayImageId": overlayImage.Id,
		"position":       map[string]float64{"x": 0.9, "y": 0.9},
		"anchor":         "BOTTOM_RIGHT",
		"opacity":        0.5,
		"scale":          0.15,
	}))
	require.NoError(t, err)
	require.NotNil(t, resp.ApplyWatermark)
	require.NotNil(t, resp.ApplyWatermark.Image)
	require.NotEmpty(t, resp.ApplyWatermark.Image.ID)
	require.Equal(t, baseImage.Name, resp.ApplyWatermark.Image.Name)
	require.Equal(t, "image/png", resp.ApplyWatermark.Image.MIMEType)

	// Verify lineage
	require.NotNil(t, resp.ApplyWatermark.Image.Parent)
	require.Equal(t, baseImage.Id, resp.ApplyWatermark.Image.Parent.ID)

	// Verify changes JSON is populated
	require.NotNil(t, resp.ApplyWatermark.Image.Changes)
	require.Contains(t, *resp.ApplyWatermark.Image.Changes, `"type"`)
	require.Contains(t, *resp.ApplyWatermark.Image.Changes, `watermark`)

	// Verify the new image exists in the DB with correct lineage
	newImage, err := tc.imageRepo.GetImageById(resp.ApplyWatermark.Image.ID)
	require.NoError(t, err)
	require.Equal(t, baseImage.Id, newImage.ParentId)
	require.Equal(t, baseImage.Id, newImage.RootId)
}

func tApplyWatermarkUnauthenticated(t *testing.T, tc *TestContext) {
	tc.clearAuthenticationInfo()

	var resp struct {
		ApplyWatermark *struct {
			Image *struct{ ID string }
		}
	}

	err := tc.client.Post(`
	mutation applyWatermark($input: ApplyWatermarkInput!) {
		applyWatermark(input: $input) {
			image {
				id
			}
		}
	}`, &resp, client.Var("input", map[string]interface{}{
		"baseImageId":    uuid.New().String(),
		"overlayImageId": uuid.New().String(),
		"position":       map[string]float64{"x": 0.5, "y": 0.5},
		"anchor":         "CENTER",
		"opacity":        1.0,
		"scale":          0.1,
	}))
	require.Error(t, err)
}

func tApplyWatermarkUnauthorizedImage(t *testing.T, tc *TestContext) {
	orgUser1 := tc.forceAuthenticate()
	orgUser2 := tc.forceAuthenticate()
	sd, tempDir := createFSStorageDefinition(t, tc)
	defer os.RemoveAll(tempDir)

	baseBytes := makeTestPNGBytes(100, 100, color.White)
	overlayBytes := makeTestPNGBytes(20, 20, color.Black)
	// base image owned by orgUser1
	baseImage := createRealImage(t, tc, orgUser1.Id, sd.Id, baseBytes)
	overlayImage := createRealImage(t, tc, orgUser2.Id, sd.Id, overlayBytes)

	// Authenticate as orgUser2 and try to edit orgUser1's image
	tc.setAuthenticatedUser(orgUser2)

	var resp struct {
		ApplyWatermark *struct {
			Image *struct{ ID string }
		}
	}

	err := tc.client.Post(`
	mutation applyWatermark($input: ApplyWatermarkInput!) {
		applyWatermark(input: $input) {
			image {
				id
			}
		}
	}`, &resp, client.Var("input", map[string]interface{}{
		"baseImageId":    baseImage.Id,
		"overlayImageId": overlayImage.Id,
		"position":       map[string]float64{"x": 0.5, "y": 0.5},
		"anchor":         "CENTER",
		"opacity":        1.0,
		"scale":          0.1,
	}))
	require.Error(t, err)
}

func tApplyWatermarkInvalidImageId(t *testing.T, tc *TestContext) {
	tc.forceAuthenticate()

	var resp struct {
		ApplyWatermark *struct {
			Image *struct{ ID string }
		}
	}

	err := tc.client.Post(`
	mutation applyWatermark($input: ApplyWatermarkInput!) {
		applyWatermark(input: $input) {
			image {
				id
			}
		}
	}`, &resp, client.Var("input", map[string]interface{}{
		"baseImageId":    uuid.New().String(),
		"overlayImageId": uuid.New().String(),
		"position":       map[string]float64{"x": 0.5, "y": 0.5},
		"anchor":         "CENTER",
		"opacity":        1.0,
		"scale":          0.1,
	}))
	require.Error(t, err)
}

func TestImageResolvers(t *testing.T) {
	tc := newTestContext(t)
	tc.runTestCases(
		tImagesNoFilterNoOrderSiteOwner,
		tSiteOwnerCanAccessAllImages,
		tNormalUserCanOnlyAcessOwnImages,
		tBasicPagination,
		tBasicPaginationByCreatedAt,
		tDeletingImage,
		tImageCreatedByIsPopulated,
		tImageCreatedByNullWhenNoCreator,
		tApplyWatermark,
		tApplyWatermarkUnauthenticated,
		tApplyWatermarkUnauthorizedImage,
		tApplyWatermarkInvalidImageId,
	)
}
