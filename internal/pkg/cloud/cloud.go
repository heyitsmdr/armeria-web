package cloud

import (
	"context"
	"io/ioutil"
	"log"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

type StorageManager struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

const (
	CharactersFile = "characters.json"
	WorldFile      = "world.json"
	MobsFile       = "mobs.json"
	ItemsFile      = "items.json"
	LedgersFile    = "ledgers.json"
)

// NewStorageManager creates a new cloud storage manager instance.
func NewStorageManager(bucketName string) *StorageManager {
	sm := &StorageManager{}

	//client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(serviceAccount)))
	client, err := storage.NewClient(context.Background(), option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create client: %v\n", err)
	}
	sm.client = client

	bucket := client.Bucket(bucketName)
	sm.bucket = bucket

	return sm
}

// ReadFile returns the contents (as []byte) of a cloud file.
func (sm *StorageManager) ReadFile(f string) []byte {
	rc, err := sm.bucket.Object(f).NewReader(context.Background())
	if err != nil {
		log.Fatalf("failed to open cloud file %s: %v", f, err)
	}
	defer rc.Close()

	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatalf("failed to read cloud file %s: %v", f, err)
	}

	return slurp
}

// WriteFile writes contents to a cloud file.
func (sm *StorageManager) WriteFile(f, ct string, contents []byte) int {
	wc := sm.bucket.Object(f).NewWriter(context.Background())
	wc.ContentType = ct
	b, err := wc.Write(contents)
	if err != nil {
		log.Fatalf("failed to write cloud file %s: %v", f, err)
		return 0
	}
	if err := wc.Close(); err != nil {
		log.Fatalf("failed to close cloud file %s: %v", f, err)
		return 0
	}
	return b
}

// CloseClient terminates the cloud client connection.
func (sm *StorageManager) CloseClient() {
	_ = sm.client.Close()
}
