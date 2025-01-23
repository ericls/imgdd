package storage_test

import (
	"testing"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/test_support"
)

func TestGetStorage(t *testing.T) {
	test_support.ResetDatabase(TestServiceMan.GetDBConfig())
	dbConn := db.GetConnection(TestServiceMan.GetDBConfig())
	repo := storage.NewDBStorageConfig(dbConn).MakeStorageDefRepo()

	storageType := "s3"
	identifier := "test"
	isEnabled := true
	priority := int64(1)
	storageDef, err := repo.CreateStorageDefinition(storageType, TestServiceMan.GetS3ConfigJSON(), identifier, isEnabled, priority)
	if err != nil {
		t.Fatal(err)
	}
	s3Store, err := storage.GetStorage(storageDef)
	if err != nil {
		t.Fatal(err)
	}
	if err := TestServiceMan.Pool.Retry(func() error {
		return s3Store.CheckConnection()
	}); err != nil {
		t.Fatal(err)
	}
}
