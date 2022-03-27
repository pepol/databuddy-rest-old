// Package kv implements simple Key-Value store on top of BadgerDB.
package kv

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/pepol/databuddy/internal/config"
)

// Store is an abstraction above single Badger database. Provides methods
// that are exported over API middleware.
type Store struct {
	name string
	path string
	db   *badger.DB
}

// Open the store with given name. The store's name is equivalent to DataBuddy
// namespace parameter in API calls.
func Open(cfg *config.Config, name string) (*Store, error) {
	if cfg.DataDir == "" {
		return nil, fmt.Errorf("data directory not set")
	}

	path := filepath.Join(cfg.DataDir, "store", name)

	opt := badger.DefaultOptions(path).
		WithCompactL0OnClose(true).
		WithMetricsEnabled(true)

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	store := new(Store)
	store.name = name
	store.path = path
	store.db = db

	return store, nil
}

// Set the given key to contain provided value, with optional (0 means disabled
// in both cases) time-to-live (TTL) and user meta byte (meta).
func (s *Store) Set(key string, value []byte, ttl time.Duration, meta byte) error {
	if s.db == nil {
		return fmt.Errorf("store '%s' not opened", s.name)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), value)

		if ttl != 0 {
			entry = entry.WithTTL(ttl)
		}

		if meta != 0 {
			entry = entry.WithMeta(meta)
		}

		if err := txn.SetEntry(entry); err != nil {
			return err
		}

		return txn.Commit()
	})
}

// Get the value currently stored at provided key.
func (s *Store) Get(key string) ([]byte, error) {
	if s.db == nil {
		return nil, fmt.Errorf("store '%s' not opened", s.name)
	}

	var value []byte

	err := s.db.View(func(txn *badger.Txn) error {
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
