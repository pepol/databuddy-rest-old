package db

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pepol/databuddy/internal/log"
)

// Database is the implementation of local storage layer.
type Database struct {
	datadir       string
	system        *Bucket
	buckets       map[string]*Bucket
	DefaultBucket string
}

const (
	bucketKeyPrefix    = "bucket:"
	datadirPermissions = 0o700
	defaultBucketKey   = "defaults:bucket"
	initKey            = "system:initialized"
	systemBucketName   = "_system"
)

// InitDatabase creates the local database for use.
func InitDatabase(datadir string, bucketName string) error {
	if err := checkDataDirectory(datadir); err != nil {
		return err
	}

	empty, err := isEmpty(datadir)
	if err != nil {
		return err
	}
	if !empty {
		return fmt.Errorf("directory '%s' not empty", datadir)
	}

	systemBucket, err := openBucketNoCheck(systemBucketName, datadir)
	if err != nil {
		return err
	}
	log.Info("created system bucket")

	if err := systemBucket.Set(bucketKeyPrefix+bucketName, []byte{1}); err != nil {
		return err
	}
	log.Info(fmt.Sprintf("created bucket '%s'", bucketName))

	if err := systemBucket.Set(defaultBucketKey, []byte(bucketName)); err != nil {
		return err
	}
	log.Info(fmt.Sprintf("set bucket '%s' as default", bucketName))

	if err := systemBucket.Set(initKey, []byte{1}); err != nil {
		return err
	}

	return systemBucket.Close()
}

// OpenDatabase opens the local database for use.
func OpenDatabase(datadir string) (*Database, error) {
	if err := checkDataDirectory(datadir); err != nil {
		return nil, err
	}

	systemBucket, err := openBucketNoCheck(systemBucketName, datadir)
	if err != nil {
		return nil, err
	}

	_, err = systemBucket.Get(initKey)
	if err != nil {
		return nil, fmt.Errorf("validating db: %v", err)
	}

	defaultBucket, err := systemBucket.Get(defaultBucketKey)
	if err != nil {
		return nil, fmt.Errorf("getting default bucket name: %v", err)
	}

	buckets := make(map[string]*Bucket)

	keys, err := systemBucket.List(bucketKeyPrefix)
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("buckets (DB): %v", keys))

	for _, key := range keys {
		if !strings.HasPrefix(key, bucketKeyPrefix) {
			continue
		}

		bucketName := strings.TrimPrefix(key, bucketKeyPrefix)

		bucket, err := openBucket(bucketName, datadir)
		if err != nil {
			log.Error(fmt.Sprintf("opening bucket '%s'", bucketName), err)
			continue
		}

		buckets[bucketName] = bucket
		log.Info(fmt.Sprintf("opened bucket '%s'", bucketName))
	}

	return &Database{
		datadir:       datadir,
		system:        systemBucket,
		buckets:       buckets,
		DefaultBucket: string(defaultBucket),
	}, nil
}

// Create creates a new database/bucket with given name.
func (db *Database) Create(name string) error {
	if err := db.system.Set(bucketKeyPrefix+name, []byte{1}); err != nil {
		return err
	}

	bucket, err := openBucket(name, db.datadir)
	if err != nil {
		return err
	}

	db.buckets[name] = bucket
	return nil
}

// Get bucket with given name, or error if it doesn't exist.
func (db *Database) Get(name string) (*Bucket, error) {
	bucket, ok := db.buckets[name]
	if !ok {
		return nil, fmt.Errorf("bucket '%s' not found", name)
	}

	return bucket, nil
}

// List all available buckets (names only).
func (db *Database) List(prefix string) []string {
	names := make([]string, 0, len(db.buckets))

	for name := range db.buckets {
		if strings.HasPrefix(name, prefix) {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names
}

// Count all available buckets.
func (db *Database) Count() int {
	return len(db.buckets)
}

// Drop given bucket.
func (db *Database) Drop(name string) error {
	if name == db.DefaultBucket {
		return fmt.Errorf("bucket '%s' is marked as default and cannot be deleted", name)
	}

	if err := db.system.Delete(bucketKeyPrefix + name); err != nil {
		return err
	}

	bucket, ok := db.buckets[name]
	delete(db.buckets, name)

	if !ok {
		return nil
	}

	return bucket.Close()
}

// Close the database.
func (db *Database) Close() error {
	for _, bucket := range db.buckets {
		if err := bucket.Close(); err != nil {
			log.Error(fmt.Sprintf("closing bucket '%s'", bucket.Name), err)
		}
	}

	return db.system.Close()
}
