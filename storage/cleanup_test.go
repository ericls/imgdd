package storage_test

import (
	"os"
	"testing"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/storage"
)

func TestCleanUp(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_fs_storage_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	configJSON := `{"mediaRoot": "` + tempDir + `"}`
	storageDefRepo := storage.NewInMemoryStorageDefRepo()
	storageDef, err := storageDefRepo.CreateStorageDefinition("fs", configJSON, "test", true, 0)
	if err != nil {
		t.FailNow()
	}
	storageInstance, err := storage.GetStorage(storageDef)
	if err != nil {
		t.FailNow()
	}
	dbConn := db.GetConnection(TestServiceMan.GetDBConfig())
	storedImageRepo := storage.NewDBStoredImageRepo(dbConn)
	imageRepo := image.NewDBImageRepo(dbConn)
	fakeImage := domainmodels.Image{
		UploaderIP:      "127.0.0.1",
		CreatedById:     "",
		MIMEType:        "image/png",
		Name:            "test.png",
		Identifier:      "test",
		NominalByteSize: int32(100),
		NominalWidth:    100,
		NominalHeight:   100,
	}
	storedImage, err := imageRepo.CreateAndSaveUploadedImage(&fakeImage, "image/png", []byte("test"), storageDef.Id, storageInstance.Save)
	if err != nil {
		t.Fatal(err)
	}
	imageId := storedImage.Image.Id
	if storedImage.IsFileDeleted {
		t.Fatal("IsFileDeleted is true")
	}
	storedImages, err := storedImageRepo.GetStoredImagesByIds([]string{storedImage.Id})
	if err != nil {
		t.Fatal(err)
	}
	storedImage = storedImages[0]
	image, _ := imageRepo.GetImageById(imageId)
	if image == nil {
		t.Fatal("image should not be nil")
	}
	err = imageRepo.DeleteImageById(imageId)
	if err != nil {
		t.Fatal(err)
	}
	image, _ = imageRepo.GetImageById(imageId)
	if image != nil {
		t.Fatal("image should be nil after deletion")
	}
	toDeleteList, err := storedImageRepo.GetStoredImagesToDelete()
	if err != nil {
		t.Fatal(err)
	}
	if len(toDeleteList) != 1 {
		t.Fatal("toDeleteList length is not 1")
	}
	count, err := storage.CleanupStoredImage(storedImageRepo, storageDefRepo)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("Deleted count is not 1")
	}
	storedImages, err = storedImageRepo.GetStoredImagesByIds([]string{storedImage.Id})
	if err != nil {
		t.Fatal(err)
	}
	storedImage = storedImages[0]
	meta := storageInstance.GetMeta(storedImage.FileIdentifier)
	if meta.ByteSize != 0 {
		t.Fatal("File is not deleted")
	}
	if !storedImage.IsFileDeleted {
		t.Fatal("IsFileDeleted is not set to true")
	}
}
