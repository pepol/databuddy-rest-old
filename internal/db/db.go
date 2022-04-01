package db

// Database is the implementation of local storage layer.
type Database struct{}

// OpenDatabase opens the local database for use.
func OpenDatabase() *Database {
	return new(Database)
}

// Create creates a new database/bucket with given name.
func (db *Database) Create(name string) error {
	// TODO: Implement bucket management.
	return nil
}

// Get bucket with given name, or error if it doesn't exist.
func (db *Database) Get(name string) (*Bucket, error) {
	// TODO: Implement bucket management.
	return nil, nil
}
