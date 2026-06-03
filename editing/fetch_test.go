package editing

import (
	"iter"
	"os"
	"path/filepath"
	"strings"
	"testing"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/google/uuid"
)

// fakeStoredImageRepo returns a fixed set of stored images for any image ID.
type fakeStoredImageRepo struct {
	images []*dm.StoredImage
}

func (r *fakeStoredImageRepo) GetStoredImagesByImageId(imageId string) ([]*dm.StoredImage, error) {
	return r.images, nil
}

func (r *fakeStoredImageRepo) GetStoredImageByIdentifierAndMimeType(identifier, mime string) ([]*dm.StoredImage, error) {
	return nil, nil
}

func (r *fakeStoredImageRepo) GetStoredImagesByIds(ids []string) ([]*dm.StoredImage, error) {
	return nil, nil
}

func (r *fakeStoredImageRepo) GetStoredImageIdsByImageIds(imageIds []string) (map[string][]string, error) {
	return nil, nil
}

func (r *fakeStoredImageRepo) GetStoredImagesToDelete() ([]*dm.StoredImage, error) {
	return nil, nil
}

func (r *fakeStoredImageRepo) MarkStoredImagesAsDeleted(ids []string) error {
	return nil
}

func (r *fakeStoredImageRepo) GetStoredImagesForReplication(sourceStorageDefId string, targetStorageDefId string) (int, iter.Seq2[*dm.StoredImage, error], error) {
	return 0, func(yield func(*dm.StoredImage, error) bool) {}, nil
}

// fakeStorageDefRepo returns a fixed set of storage definitions by ID.
type fakeStorageDefRepo struct {
	defs map[string]*dm.StorageDefinition
}

func (r *fakeStorageDefRepo) GetStorageDefinitionsByIds(ids []string) ([]*dm.StorageDefinition, error) {
	result := make([]*dm.StorageDefinition, 0, len(ids))
	for _, id := range ids {
		if def, ok := r.defs[id]; ok {
			result = append(result, def)
		}
	}
	return result, nil
}

func (r *fakeStorageDefRepo) GetStorageDefinitionById(id string) (*dm.StorageDefinition, error) {
	return r.defs[id], nil
}

func (r *fakeStorageDefRepo) GetStorageDefinitionByIdentifier(identifier string) (*dm.StorageDefinition, error) {
	for _, def := range r.defs {
		if def.Identifier == identifier {
			return def, nil
		}
	}
	return nil, nil
}

func (r *fakeStorageDefRepo) ListStorageDefinitions() ([]*dm.StorageDefinition, error) {
	result := make([]*dm.StorageDefinition, 0, len(r.defs))
	for _, def := range r.defs {
		result = append(result, def)
	}
	return result, nil
}

func (r *fakeStorageDefRepo) CreateStorageDefinition(storage_type string, config string, identifier string, isEnabled bool, priority int64) (*dm.StorageDefinition, error) {
	return nil, nil
}

func (r *fakeStorageDefRepo) UpdateStorageDefinition(identifier string, storage_type *string, config *string, isEnabled *bool, priority *int64) (*dm.StorageDefinition, error) {
	return nil, nil
}

func makeFetchTestRepos(t *testing.T, fileContent []byte) (*fakeStoredImageRepo, *fakeStorageDefRepo, string) {
	t.Helper()
	mediaRoot := t.TempDir()
	fileIdentifier := "test-image.bin"
	if err := os.WriteFile(filepath.Join(mediaRoot, fileIdentifier), fileContent, 0644); err != nil {
		t.Fatalf("failed to write test image file: %v", err)
	}

	defID := uuid.New().String()
	storageDef := &dm.StorageDefinition{
		Id:          defID,
		Identifier:  "test-fs",
		StorageType: dm.FSStorageType,
		Config:      `{"mediaRoot":"` + mediaRoot + `"}`,
		IsEnabled:   true,
		Priority:    0,
	}
	storedImage := &dm.StoredImage{
		Id:                  uuid.New().String(),
		FileIdentifier:      fileIdentifier,
		StorageDefinitionId: defID,
		IsFileDeleted:       false,
	}

	storedImageRepo := &fakeStoredImageRepo{images: []*dm.StoredImage{storedImage}}
	storageDefRepo := &fakeStorageDefRepo{defs: map[string]*dm.StorageDefinition{defID: storageDef}}
	return storedImageRepo, storageDefRepo, uuid.New().String()
}

func TestFetchImageFuncRejectsOversizedImage(t *testing.T) {
	content := make([]byte, 200)
	storedImageRepo, storageDefRepo, imageID := makeFetchTestRepos(t, content)

	fetch := NewFetchImageFunc(storedImageRepo, storageDefRepo, 100)
	_, err := fetch(imageID)
	if err == nil {
		t.Fatal("expected error for image exceeding maxBytes, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Fatalf("expected 'exceeds maximum size' in error, got: %v", err)
	}
}

func TestFetchImageFuncAcceptsImageWithinLimit(t *testing.T) {
	content := make([]byte, 100)
	storedImageRepo, storageDefRepo, imageID := makeFetchTestRepos(t, content)

	fetch := NewFetchImageFunc(storedImageRepo, storageDefRepo, 200)
	data, err := fetch(imageID)
	if err != nil {
		t.Fatalf("expected no error for image within maxBytes, got: %v", err)
	}
	if len(data) != len(content) {
		t.Fatalf("expected %d bytes, got %d", len(content), len(data))
	}
}
