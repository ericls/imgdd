package storage_test

import (
	"os"
	"testing"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/test_support"
)

func TestBulkReplicateToStorageDefinition(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "test_replication_src_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "test_replication_dst_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	test_support.ResetDatabase(TestServiceMan.GetDBConfig())
	dbConn := db.GetConnection(TestServiceMan.GetDBConfig())
	storageDefRepo := storage.NewInMemoryStorageDefRepo()
	storedImageRepo := storage.NewDBStoredImageRepo(dbConn)
	imageRepo := image.NewDBImageRepo(dbConn)

	srcDef, err := storageDefRepo.CreateStorageDefinition("fs", `{"mediaRoot":"`+srcDir+`"}`, "src", true, 0)
	if err != nil {
		t.Fatal(err)
	}
	dstDef, err := storageDefRepo.CreateStorageDefinition("fs", `{"mediaRoot":"`+dstDir+`"}`, "dst", true, 1)
	if err != nil {
		t.Fatal(err)
	}

	srcStorage, err := storage.GetStorage(srcDef)
	if err != nil {
		t.Fatal(err)
	}

	fakeImage := domainmodels.Image{
		UploaderIP:      "127.0.0.1",
		MIMEType:        "image/png",
		Name:            "test.png",
		Identifier:      "testreplica",
		NominalByteSize: 4,
		NominalWidth:    1,
		NominalHeight:   1,
	}
	storedImage, err := imageRepo.CreateAndSaveUploadedImage(&fakeImage, "image/png", []byte("test"), srcDef.Id, srcStorage.Save)
	if err != nil {
		t.Fatal(err)
	}

	count, err := storage.BulkReplicateToStorageDefinition(srcDef.Id, dstDef.Id, storedImageRepo, imageRepo, storageDefRepo, 10)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 replicated image, got %d", count)
	}

	// More directly: verify there's a stored image record for the image in the dst backend
	imageStoredImages, err := storedImageRepo.GetStoredImagesByImageId(storedImage.Image.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(imageStoredImages) != 2 {
		t.Fatalf("expected 2 stored images (src + dst), got %d", len(imageStoredImages))
	}

	var dstStoredImage *domainmodels.StoredImage
	for _, si := range imageStoredImages {
		if si.StorageDefinitionId == dstDef.Id {
			dstStoredImage = si
		}
	}
	if dstStoredImage == nil {
		t.Fatal("no stored image found in destination backend")
	}

	// Verify the file actually exists on the destination storage
	dstStorage, err := storage.GetStorage(dstDef)
	if err != nil {
		t.Fatal(err)
	}
	meta := dstStorage.GetMeta(dstStoredImage.FileIdentifier)
	if meta.ByteSize == 0 {
		t.Fatal("replicated file has zero byte size in destination storage")
	}

	// Running again should replicate 0 (already exists)
	count, err = storage.BulkReplicateToStorageDefinition(srcDef.Id, dstDef.Id, storedImageRepo, imageRepo, storageDefRepo, 10)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0 on second run (already replicated), got %d", count)
	}
}

func TestReplicateImageToStorageDefinition(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "test_replication_src_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "test_replication_dst_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	test_support.ResetDatabase(TestServiceMan.GetDBConfig())
	dbConn := db.GetConnection(TestServiceMan.GetDBConfig())
	storageDefRepo := storage.NewInMemoryStorageDefRepo()
	storedImageRepo := storage.NewDBStoredImageRepo(dbConn)
	imageRepo := image.NewDBImageRepo(dbConn)

	srcDef, err := storageDefRepo.CreateStorageDefinition("fs", `{"mediaRoot":"`+srcDir+`"}`, "src2", true, 0)
	if err != nil {
		t.Fatal(err)
	}
	dstDef, err := storageDefRepo.CreateStorageDefinition("fs", `{"mediaRoot":"`+dstDir+`"}`, "dst2", true, 1)
	if err != nil {
		t.Fatal(err)
	}

	srcStorage, err := storage.GetStorage(srcDef)
	if err != nil {
		t.Fatal(err)
	}

	fakeImage := domainmodels.Image{
		UploaderIP:      "127.0.0.1",
		MIMEType:        "image/png",
		Name:            "test2.png",
		Identifier:      "testreplica2",
		NominalByteSize: 4,
		NominalWidth:    1,
		NominalHeight:   1,
	}
	storedImage, err := imageRepo.CreateAndSaveUploadedImage(&fakeImage, "image/png", []byte("test"), srcDef.Id, srcStorage.Save)
	if err != nil {
		t.Fatal(err)
	}
	imageId := storedImage.Image.Id

	newSI, err := storage.ReplicateImageToStorageDefinition(imageId, srcDef.Id, dstDef.Id, storedImageRepo, imageRepo, storageDefRepo)
	if err != nil {
		t.Fatal(err)
	}
	if newSI.StorageDefinitionId != dstDef.Id {
		t.Fatalf("expected stored image in dst, got storage def %s", newSI.StorageDefinitionId)
	}

	// Replicating again should fail
	_, err = storage.ReplicateImageToStorageDefinition(imageId, srcDef.Id, dstDef.Id, storedImageRepo, imageRepo, storageDefRepo)
	if err == nil {
		t.Fatal("expected error when replicating to backend that already has the image")
	}
}
