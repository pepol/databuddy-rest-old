package db

import "sync"

// DefaultBucketName contains the name of bucket created on database initialization.
const DefaultBucketName = "default"

// Bucket is the single "table" within the database.
type Bucket struct {
	Name string

	Mutex sync.RWMutex
}

// Get value stored under key.
func (b *Bucket) Get(key string) ([]byte, error) {
	// TODO: Implement.
	return nil, nil
}

// Set key to point to value.
func (b *Bucket) Set(key string, value []byte) error {
	// TODO: Implement.
	return nil
}

// Delete value stored under key.
func (b *Bucket) Delete(key string) error {
	// TODO: Implement.
	return nil
}
