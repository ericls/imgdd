package storage_test

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/storage"
)

type noopLock struct{}

func (noopLock) AcquireLock() (bool, error) { return true, nil }
func (noopLock) ReleaseLock() error         { return nil }

func TestImageExpiry(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_fs_storage_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	storageDefRepo := storage.NewInMemoryStorageDefRepo()
	storageDef, err := storageDefRepo.CreateStorageDefinition("fs", `{"mediaRoot": "`+tempDir+`"}`, "expiry", true, 0)
	if err != nil {
		t.Fatal(err)
	}
	storageInstance, err := storage.GetStorage(storageDef)
	if err != nil {
		t.Fatal(err)
	}
	dbConn := db.GetConnection(TestServiceMan.GetDBConfig())
	storedImageRepo := storage.NewDBStoredImageRepo(dbConn)
	imageRepo := image.NewDBImageRepo(dbConn)

	upload := func(identifier string, expiresAt *time.Time) *domainmodels.StoredImage {
		t.Helper()
		img := domainmodels.Image{
			MIMEType:   "image/png",
			Name:       identifier + ".png",
			Identifier: identifier,
			ExpiresAt:  expiresAt,
		}
		si, err := imageRepo.CreateAndSaveUploadedImage(&img, "image/png", []byte("test"), storageDef.Id, storageInstance.Save)
		if err != nil {
			t.Fatal(err)
		}
		return si
	}
	isServed := func(identifier string) bool {
		t.Helper()
		sis, err := storedImageRepo.GetStoredImageByIdentifierAndMimeType(identifier, "image/png")
		if err != nil {
			t.Fatal(err)
		}
		return len(sis) > 0
	}

	future := time.Now().Add(time.Hour).Truncate(time.Microsecond)
	expiring := upload("expiry-expiring", &future)
	permanent := upload("expiry-permanent", nil)

	if expiring.Image.ExpiresAt == nil || !expiring.Image.ExpiresAt.Equal(future) {
		t.Fatalf("ExpiresAt not persisted: got %v, want %v", expiring.Image.ExpiresAt, future)
	}
	if permanent.Image.ExpiresAt != nil {
		t.Fatalf("expected no expiry, got %v", permanent.Image.ExpiresAt)
	}
	if !isServed("expiry-expiring") {
		t.Fatal("image with future expiry should be served")
	}

	// Clearing and re-setting the expiry.
	if err := imageRepo.SetImageExpiration(expiring.Image.Id, nil); err != nil {
		t.Fatal(err)
	}
	img, _ := imageRepo.GetImageById(expiring.Image.Id)
	if img == nil || img.ExpiresAt != nil {
		t.Fatalf("expiry should be cleared, got %+v", img)
	}
	past := time.Now().Add(-time.Minute)
	if err := imageRepo.SetImageExpiration(expiring.Image.Id, &past); err != nil {
		t.Fatal(err)
	}

	// Expired images are hidden immediately, before any sweep runs.
	if img, _ := imageRepo.GetImageById(expiring.Image.Id); img != nil {
		t.Fatal("expired image should not be returned by GetImageById")
	}
	if isServed("expiry-expiring") {
		t.Fatal("expired image should not be served")
	}
	if !isServed("expiry-permanent") {
		t.Fatal("non-expiring image should still be served")
	}
	if err := imageRepo.SetImageExpiration(expiring.Image.Id, &future); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows when extending an expired image, got %v", err)
	}

	// The cleanup task soft-deletes the expired image and removes its file.
	if err := storage.CleanupStoredImageTask(noopLock{}, imageRepo, storedImageRepo, storageDefRepo); err != nil {
		t.Fatal(err)
	}
	sis, err := storedImageRepo.GetStoredImagesByIds([]string{expiring.Id, permanent.Id})
	if err != nil {
		t.Fatal(err)
	}
	for _, si := range sis {
		switch si.Id {
		case expiring.Id:
			if !si.IsFileDeleted {
				t.Fatal("expired image's file should be marked deleted")
			}
			if storageInstance.GetMeta(si.FileIdentifier).ByteSize != 0 {
				t.Fatal("expired image's file should be removed from storage")
			}
		case permanent.Id:
			if si.IsFileDeleted {
				t.Fatal("non-expiring image's file should be kept")
			}
		}
	}

	count, err := imageRepo.DeleteExpiredImages()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("already-swept images should not be deleted again, got %d", count)
	}
}
