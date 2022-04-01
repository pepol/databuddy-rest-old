// Package server contains the main server code, including RESP handling
// and high-level database operations.
package server

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/pepol/databuddy/internal/context"
	"github.com/pepol/databuddy/internal/db"
	"github.com/pepol/databuddy/internal/log"
	"github.com/spf13/viper"
	"github.com/tidwall/redcon"
)

type commandInfo struct {
	usage string
	help  string
}

// Handler is the main server connection handler.
type Handler struct {
	mutex               sync.RWMutex
	db                  *db.Database
	commandDescriptions map[string]commandInfo

	Mux *redcon.ServeMux

	addr     string
	hostname string
	version  string
}

// NewHandler initialized the server Handler.
func NewHandler(version, addr, hostname string) *Handler {
	return &Handler{
		commandDescriptions: make(map[string]commandInfo),
		db:                  db.OpenDatabase(),
		Mux:                 redcon.NewServeMux(),
		addr:                addr,
		hostname:            hostname,
		version:             version,
	}
}

// Serve the database over network.
func Serve(version string) {
	port := viper.GetInt("port")
	host := viper.GetString("host")

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	hostname, err := os.Hostname()
	if err != nil {
		log.Error("getting hostname", err)
		hostname = "localhost"
	}

	handler := NewHandler(version, addr, hostname)

	// General commands.
	handler.RegisterCommand("info", handler.info, "INFO [<command> ...]", "show information about command(s)")
	handler.RegisterCommand("node", handler.nodeInfo, "NODE", "return information about current node")
	handler.RegisterCommand("ping", handler.ping, "PING", "respond with 'PONG'")
	handler.RegisterCommand("quit", handler.quit, "QUIT", "close the connection")

	// DB management commands.
	handler.RegisterCommand("create", handler.create, "CREATE <database>", "create database with given name")
	handler.RegisterCommand("use", handler.use, "USE <database>", "set database for further queries")

	// KV commands.
	handler.RegisterCommand("get", handler.get, "GET <key>", "return value stored under given key")
	handler.RegisterCommand("set", handler.set, "SET <key> <value>", "store value under key, returns 'OK' if successful, 'ERR' otherwise")
	handler.RegisterCommand("del", handler.del, "DEL <key> [<key> ...]", "delete values stored under key(s), returns number of deleted items")

	log.Info(fmt.Sprintf("Starting DataBuddy %s RESP server on %s", version, addr))

	err = redcon.ListenAndServe(
		addr,
		handler.Mux.ServeRESP,
		handler.acceptConnection,
		func(conn redcon.Conn, err error) {},
	)
	if err != nil {
		log.Error("serving resp", err)
	}
}

// RegisterCommand registers command into RESP handler with given handler,
// usage information, and more detailed help text.
func (h *Handler) RegisterCommand(command string, handler redcon.HandlerFunc, usage string, help string) *Handler {
	h.commandDescriptions[command] = commandInfo{
		usage: usage,
		help:  help,
	}
	h.Mux.HandleFunc(command, handler)

	return h
}

// Initialize connection context on connection accept.
func (h *Handler) acceptConnection(conn redcon.Conn) bool {
	bucket, err := h.db.Get(db.DefaultBucketName)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR initializing connection: %v", err))
		return false
	}

	conn.SetContext(&context.Context{
		Bucket: bucket,
	})

	return true
}

// INFO [<command> ...]
// Show information about given commands (or list all of them).
func (h *Handler) info(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) == 1 {
		conn.WriteArray(len(h.commandDescriptions) + 1)
		for _, info := range h.commandDescriptions {
			conn.WriteBulkString(info.usage)
		}
		conn.WriteBulkString("\r\n")
		return
	}

	conn.WriteArray(len(cmd.Args)) // This should be 'len - 1', but we're adding additional newline string.
	for _, argB := range cmd.Args[1:] {
		arg := string(argB)
		info, ok := h.commandDescriptions[arg]
		if !ok {
			conn.WriteError(fmt.Sprintf("ERR command not found '%s'", arg))
			continue
		}
		conn.WriteBulkString(fmt.Sprintf("%s\r\n\t%s", info.usage, info.help))
	}
	conn.WriteBulkString("\r\n")
}
