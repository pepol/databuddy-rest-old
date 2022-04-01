// Package context contains connection-specific information.
package context

import "github.com/pepol/databuddy/internal/db"

// Context contains information pertaining to the current connection.
type Context struct {
	Bucket *db.Bucket
}
