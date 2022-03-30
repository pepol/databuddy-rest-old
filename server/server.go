// Package server contains the main server code, including RESP handling
// and high-level database operations.
package server

import (
	"sync"

	"github.com/pepol/databuddy/internal/db"
	"github.com/pepol/databuddy/internal/log"
	"github.com/tidwall/redcon"
)

// Handler is the main server connection handler.
type Handler struct {
	mutex sync.RWMutex
	db    *db.Database
}

// NewHandler initialized the server Handler.
func NewHandler() *Handler {
	return &Handler{
		db: db.OpenDatabase(),
	}
}

var addr = ":6543"

// Serve the database over network.
func Serve() {
	handler := NewHandler()
	mux := redcon.NewServeMux()

	// General commands.
	mux.HandleFunc("ping", handler.ping)
	mux.HandleFunc("quit", handler.quit)

	// KV commands.
	mux.HandleFunc("get", handler.get)
	mux.HandleFunc("set", handler.set)
	mux.HandleFunc("del", handler.del)

	err := redcon.ListenAndServe(
		addr,
		mux.ServeRESP,
		func(conn redcon.Conn) bool {
			return true
		},
		func(conn redcon.Conn, err error) {},
	)
	if err != nil {
		log.Error("serving resp", err)
	}
}
