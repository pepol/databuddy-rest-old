package db

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

// DefaultBucketName contains the name of bucket created on database initialization.
const DefaultBucketName = "default"

// Bucket is the single "table" within the database.
type Bucket struct {
	Name string

	path  string
	db    *badger.DB
	mutex sync.RWMutex
}

const (
	rfc1123LabelRegexFmt  = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
	rfc1123LabelMaxLength = 63
)

func openBucket(name, basePath string) (*Bucket, error) {
	if basePath == "" {
		return nil, fmt.Errorf("no path specified for bucket %s", name)
	}

	if !isValidBucketName(name) {
		return nil, fmt.Errorf("bucket name '%s' does not match RFC1123 label requirements", name)
	}

	return openBucketNoCheck(name, basePath)
}

func openBucketNoCheck(name, basePath string) (*Bucket, error) {
	path, err := filepath.Abs(filepath.Join(basePath, "buckets", name))
	if err != nil {
		return nil, err
	}

	opt := badger.DefaultOptions(path).
		WithCompactL0OnClose(true).
		WithMetricsEnabled(true)

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		Name: name,
		path: path,
		db:   db,
	}, nil
}

// List keys with given prefix.
func (b *Bucket) List(prefix string) ([]string, error) {
	if b.db == nil {
		return nil, fmt.Errorf("bucket '%s' not opened", b.Name)
	}

	var buckets []string

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		prefixB := []byte(prefix)

		for it.Seek(prefixB); it.ValidForPrefix(prefixB); it.Next() {
			item := it.Item()
			key := string(item.KeyCopy(nil))
			buckets = append(buckets, key)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return buckets, nil
}

// Get value stored under key.
func (b *Bucket) Get(key string) ([]byte, error) {
	if b.db == nil {
		return nil, fmt.Errorf("bucket '%s' not opened", b.Name)
	}

	var value []byte

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Set key to point to value.
func (b *Bucket) Set(key string, value []byte) error {
	if b.db == nil {
		return fmt.Errorf("bucket '%s' not opened", b.Name)
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), value)

		if err := txn.SetEntry(entry); err != nil {
			return err
		}

		return nil
	})
}

// Delete value stored under key.
func (b *Bucket) Delete(key string) error {
	if b.db == nil {
		return fmt.Errorf("bucket '%s' not opened", b.Name)
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(key)); err != nil {
			return err
		}

		return txn.Commit()
	})
}

// Close the underlying BadgerDB.
func (b *Bucket) Close() error {
	if b.db == nil {
		return nil // Database isn't even opened.
	}

	return b.db.Close()
}

func isValidBucketName(name string) bool {
	rfc1123LabelRegex := regexp.MustCompile("^" + rfc1123LabelRegexFmt + "$")

	if len(name) > rfc1123LabelMaxLength {
		return false
	}

	if !rfc1123LabelRegex.MatchString(name) {
		return false
	}

	return true
}
