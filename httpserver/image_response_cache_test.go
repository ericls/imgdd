package httpserver

import (
	"iter"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/storage"
)

type testStoredImageRepo struct {
	count       int
	storedImage *domainmodels.StoredImage
}

func (r *testStoredImageRepo) GetStoredImageByIdentifierAndMimeType(identifier, mime string) ([]*domainmodels.StoredImage, error) {
	r.count++
	return []*domainmodels.StoredImage{r.storedImage}, nil
}

func (r *testStoredImageRepo) GetStoredImagesByIds(ids []string) ([]*domainmodels.StoredImage, error) {
	return nil, nil
}

func (r *testStoredImageRepo) GetStoredImageIdsByImageIds(imageIds []string) (map[string][]string, error) {
	return nil, nil
}

func (r *testStoredImageRepo) GetStoredImagesToDelete() ([]*domainmodels.StoredImage, error) {
	return nil, nil
}

func (r *testStoredImageRepo) MarkStoredImagesAsDeleted(ids []string) error {
	return nil
}

func (r *testStoredImageRepo) GetStoredImagesByImageId(imageId string) ([]*domainmodels.StoredImage, error) {
	return nil, nil
}

func (r *testStoredImageRepo) GetStoredImagesForReplication(sourceStorageDefId string, targetStorageDefId string) (int, iter.Seq2[*domainmodels.StoredImage, error], error) {
	return 0, func(yield func(*domainmodels.StoredImage, error) bool) {}, nil
}

func TestImageHandlerUsesConfiguredResponseCache(t *testing.T) {
	mediaRoot := t.TempDir()
	fileIdentifier := "stored.png"
	if err := os.WriteFile(filepath.Join(mediaRoot, fileIdentifier), []byte("first"), 0644); err != nil {
		t.Fatalf("failed to write image: %v", err)
	}

	storageDef := &domainmodels.StorageDefinition{
		Id:          "00000000-0000-0000-0000-000000000001",
		Identifier:  "fs1",
		StorageType: domainmodels.FSStorageType,
		Config:      `{"mediaRoot":"` + mediaRoot + `"}`,
		IsEnabled:   true,
		Priority:    0,
	}
	storageDefRepo := storage.NewInMemoryStorageDefRepo()
	storageDefRepo.AddStorageDefinition(storageDef)
	storedImageRepo := &testStoredImageRepo{
		storedImage: &domainmodels.StoredImage{
			Id:                  "00000000-0000-0000-0000-000000000002",
			FileIdentifier:      fileIdentifier,
			StorageDefinitionId: storageDef.Id,
			IsFileDeleted:       false,
		},
	}
	cache := newImageResponseCache(1024, 1024)
	handler := makeImageHandler(storageDefRepo, storedImageRepo, cache)

	req := httptest.NewRequest(http.MethodGet, "/image/hot.png", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected first request status 200, got %d", res.Code)
	}
	if body := res.Body.String(); body != "first" {
		t.Fatalf("expected first body from storage, got %q", body)
	}

	if err := os.WriteFile(filepath.Join(mediaRoot, fileIdentifier), []byte("later"), 0644); err != nil {
		t.Fatalf("failed to rewrite image: %v", err)
	}
	res = httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected second request status 200, got %d", res.Code)
	}
	if body := res.Body.String(); body != "first" {
		t.Fatalf("expected second body from response cache, got %q", body)
	}
	if storedImageRepo.count != 1 {
		t.Fatalf("expected cache hit to skip stored image lookup, got %d lookups", storedImageRepo.count)
	}
}

func TestImageResponseCacheEvictsByByteBudget(t *testing.T) {
	cache := newImageResponseCache(6, 6)
	cache.put("a", "text/plain", "a", "fs1", []byte("1234"))
	cache.put("b", "text/plain", "b", "fs1", []byte("5678"))

	if _, ok := cache.get("a"); ok {
		t.Fatal("expected oldest entry to be evicted")
	}
	if entry, ok := cache.get("b"); !ok || string(entry.body) != "5678" {
		t.Fatal("expected newest entry to remain")
	}
}

func TestImageResponseCacheSkipsOversizedFiles(t *testing.T) {
	cache := newImageResponseCache(10, 3)
	cache.put("a", "text/plain", "a", "fs1", []byte("1234"))

	if _, ok := cache.get("a"); ok {
		t.Fatal("expected oversized entry not to be cached")
	}
}

var _ storage.StoredImageRepo = (*testStoredImageRepo)(nil)
