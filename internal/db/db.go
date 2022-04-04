package db

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pepol/databuddy/internal/log"
)

// Database is the implementation of local storage layer.
type Database struct {
	datadir string
	system  *Bucket
	buckets map[string]*Bucket
}

const (
	bucketKeyPrefix    = "bucket:"
	datadirPermissions = 0o700
	systemBucketName   = "_system"
)

// OpenDatabase opens the local database for use.
func OpenDatabase(datadir string) (*Database, error) {
	if err := checkDataDirectory(datadir); err != nil {
		return nil, err
	}

	systemBucket, err := openBucketNoCheck(systemBucketName, datadir)
	if err != nil {
		return nil, err
	}

	buckets := make(map[string]*Bucket)

	keys, err := systemBucket.List(bucketKeyPrefix)
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("buckets (DB): %v", keys))

	// TODO: Move this into "init database" subcommand.
	// Create default bucket if no bucket exists.
	if len(keys) == 0 {
		if err := systemBucket.Set(bucketKeyPrefix+DefaultBucketName, []byte{1}); err != nil {
			return nil, err
		}
		keys = append(keys, bucketKeyPrefix+DefaultBucketName)
		log.Info(fmt.Sprintf("created default bucket '%s'", DefaultBucketName))
	}

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
		datadir: datadir,
		system:  systemBucket,
		buckets: buckets,
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

func checkDataDirectory(datadir string) error {
	fi, err := os.Stat(datadir)

	if os.IsNotExist(err) {
		if err = os.Mkdir(datadir, datadirPermissions); err != nil {
			return err
		}
		fi, err = os.Stat(datadir)
	}
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("'%s' exists and is not a directory", datadir)
	}

	if fi.Mode().Perm() != datadirPermissions {
		return fmt.Errorf("permissions for '%s' are incorrect (%o != %o)", datadir, fi.Mode(), datadirPermissions)
	}

	return nil
}
