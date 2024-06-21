package storage_test

import (
	"fmt"
	"imgdd/db"
	"imgdd/storage"
	"imgdd/test_support"
	"testing"
)

func TestRepoCRUD(t *testing.T) {
	test_support.ResetDatabase(&TEST_DB_CONF)

	// Create a new storage repo
	dbConn := db.GetConnection(&TEST_DB_CONF)
	repo := storage.NewDBStorageRepo(dbConn)

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
	config := fmt.Sprintf(`{"endpoint":"http://localhost:%s","bucket":"%s","access":"%s","secret":"%s"}`,
		testS3Port, testS3Bucket, testS3Access, testS3Secret,
	)
	identifier := "test"
	isEnabled := true
	priority := int64(1)
	storageDef, err := repo.CreateStorageDefinition(storageType, config, identifier, isEnabled, priority)
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
