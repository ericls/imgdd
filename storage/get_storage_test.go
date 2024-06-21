package storage_test

import (
	"fmt"
	"imgdd/db"
	"imgdd/storage"
	"imgdd/test_support"
	"testing"
)

func TestGetStorage(t *testing.T) {
	test_support.ResetDatabase(&TEST_DB_CONF)
	dbConn := db.GetConnection(&TEST_DB_CONF)
	repo := storage.NewDBStorageRepo(dbConn)

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
	s3Store, err := storage.GetStorage(storageDef)
	if err != nil {
		t.Fatal(err)
	}
	if err := dockerTestPool.Retry(func() error {
		return s3Store.CheckConnection()
	}); err != nil {
		t.Fatal(err)
	}
}
