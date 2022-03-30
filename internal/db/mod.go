// Package db implements high-level database driver (multiple badger "buckets").
package db

// Database is the implementation of local storage layer.
type Database struct{}

// OpenDatabase opens the local database for use.
func OpenDatabase() *Database {
	return new(Database)
}

// Get a value stored under key.
func (db *Database) Get(_key string) ([]byte, error) {
	return nil, nil
}

// Set a key to contain the value.
func (db *Database) Set(_key string, _value []byte) error {
	return nil
}

// Delete the value stored under key.
func (db *Database) Delete(_key string) error {
	return nil
}
