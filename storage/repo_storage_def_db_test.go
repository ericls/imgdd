package storage_test

import (
	"testing"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/test_support"
)

func runCRUDTest(t *testing.T, repo storage.StorageDefRepo) {
	// List storage definitions
	storageDefs, err := repo.ListStorageDefinitions()
	if err != nil {
		t.Fatal(err)
	}
	if len(storageDefs) != 0 {
		t.Fatal("storage definitions not empty. (Sanity check)")
	}

	// Create a new storage definition
	storageType := "s3"
	identifier := "test"
	isEnabled := true
	priority := int64(1)
	storageDef, err := repo.CreateStorageDefinition(storageType, TestServiceMan.GetS3ConfigJSON(), identifier, isEnabled, priority)
	if err != nil {
		t.Fatal(err)
	}
	if storageDef == nil {
		t.Fatal("storage definition not found")
	}

	// Get storage definition by ID
	storageDefByID, err := repo.GetStorageDefinitionById(storageDef.Id)
	if err != nil {
		t.Fatal(err)
	}
	if storageDefByID == nil {
		t.Fatal("storage definition by ID not found")
	}

	// Get storage definition by identifier
	storageDefByIdentifier, err := repo.GetStorageDefinitionByIdentifier(storageDef.Identifier)
	if err != nil {
		t.Fatal(err)
	}
	if storageDefByIdentifier == nil {
		t.Fatal("storage definition by identifier not found")
	}

	// Update storage definition
	newStorageType := "s3"
	newConfig := `{"endpoint":"http://localhost:9000","bucket":"test","access":"minio","secret":"minio123"}`
	newIsEnabled := false
	newPriority := int64(2)
	updatedStorageDef, err := repo.UpdateStorageDefinition(storageDef.Identifier, &newStorageType, &newConfig, &newIsEnabled, &newPriority)
	if err != nil {
		t.Fatal(err)
	}
	if updatedStorageDef == nil {
		t.Fatal("updated storage definition not found")
	}
}

func TestRepoCRUD(t *testing.T) {
	test_support.ResetDatabase(TestServiceMan.GetDBConfig())

	// Create a new storage repo
	dbConn := db.GetConnection(TestServiceMan.GetDBConfig())
	repoDb := storage.NewDBStorageDefRepo(dbConn)
	runCRUDTest(t, repoDb)
	repoInMemory := storage.NewInMemoryStorageDefRepo()
	runCRUDTest(t, repoInMemory)
}
