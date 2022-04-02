// Package server contains the main server code, including RESP handling
// and high-level database operations.
package server

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/pepol/databuddy/internal/context"
	"github.com/pepol/databuddy/internal/db"
	"github.com/pepol/databuddy/internal/log"
	"github.com/spf13/viper"
	"github.com/tidwall/redcon"
)

type commandInfo struct {
	arity      int
	flags      []string
	firstKey   int
	lastKey    int
	stepKey    int
	categories []string
	tips       []string
}

// Handler is the main server connection handler.
type Handler struct {
	mutex sync.RWMutex
	db    *db.Database

	// Replace with sorted map implementation for consistent ordering.
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
//nolint:gomnd // Magic numbers in command registration calls are not magic.
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

	// Meta (command-handling) commands.
	handler.Register("command", handler.command, 1, []string{"server"}, -1, -1, 0, nil, []string{"COMMAND", "show information about all commands"})
	handler.RegisterChild("command count", 2, []string{"server"}, -1, -1, 0, nil, []string{"COMMAND COUNT", "return count of all available commands"})
	handler.RegisterChild("command list", 2, []string{"server"}, -1, -1, 0, nil, []string{"COMMAND LIST", "return list of all available commands"})
	handler.RegisterChild("command info", -2, []string{"server"}, 2, -1, 1, nil, []string{"COMMAND INFO [<command> ...]", "show system information about given command(s)"})
	handler.RegisterChild("command docs", -2, []string{"server"}, 2, -1, 1, nil, []string{"COMMAND DOCS [<command> ...]", "show documentation about given command(s)"})

	// General information commands.
	handler.Register("info", handler.info, 1, []string{"server"}, -1, -1, 0, nil, []string{"INFO", "return information about current node"})
	handler.Register("ping", handler.ping, 1, []string{"general"}, -1, -1, 0, nil, []string{"PING", "respond with 'PONG'"})
	handler.Register("quit", handler.quit, 1, []string{"server"}, -1, -1, 0, nil, []string{"QUIT", "close the connection"})

	// DB management commands.
	handler.Register("create", handler.create, 2, []string{"database"}, 1, 1, 0, nil, []string{"CREATE <database>", "create database with given name"})
	handler.Register("use", handler.use, 2, []string{"database"}, 1, 1, 0, nil, []string{"USE <database>", "set database for further queries"})

	// KV commands.
	handler.Register("get", handler.get, 2, []string{"read"}, 1, 1, 0, nil, []string{"GET <key>", "return value stored under given key"})
	handler.Register("set", handler.set, 3, []string{"write"}, 1, 1, 0, nil, []string{"SET <key> <value>", "store value under key, returns 'OK' if successful, 'ERR' otherwise"})
	handler.Register("del", handler.del, -2, []string{"write"}, 1, -1, 1, nil, []string{"DEL <key> [<key> ...]", "delete values stored under key(s), returns number of deleted items"})

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

// Register command into RESP handler with given handler and usage information
// for the command.
func (h *Handler) Register(
	command string,
	handler redcon.HandlerFunc,
	arity int,
	flags []string,
	firstKey int,
	lastKey int,
	stepKey int,
	categories []string,
	tips []string,
) *Handler {
	h.commandDescriptions[command] = commandInfo{
		arity:      arity,
		flags:      flags,
		firstKey:   firstKey,
		lastKey:    lastKey,
		stepKey:    stepKey,
		categories: categories,
		tips:       tips,
	}
	h.Mux.HandleFunc(command, handler)

	return h
}

// RegisterChild registers sub-command's usage information only.
// Handler for parent command needs to handle the sub-command call!
func (h *Handler) RegisterChild(
	command string,
	arity int,
	flags []string,
	firstKey int,
	lastKey int,
	stepKey int,
	categories []string,
	tips []string,
) *Handler {
	h.commandDescriptions[command] = commandInfo{
		arity:      arity,
		flags:      flags,
		firstKey:   firstKey,
		lastKey:    lastKey,
		stepKey:    stepKey,
		categories: categories,
		tips:       tips,
	}

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

// COMMAND [<command> ...]
// Show information about given commands (or list all of them).
func (h *Handler) command(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) == 1 {
		h.commandInfo(conn, nil)
		return
	}

	subcommand := strings.ToLower(string(cmd.Args[1]))

	switch subcommand {
	case "count":
		h.commandCount(conn)
	case "list":
		h.commandList(conn)
	case "info":
		h.commandInfo(conn, cmd.Args[2:])
	case "docs":
		h.commandDocs(conn, cmd.Args[2:])
	default:
		conn.WriteError(fmt.Sprintf("ERR unknown command '%s %s'", string(cmd.Args[0]), subcommand))
	}
}

// COMMAND COUNT
// Return count of all commands available on server.
func (h *Handler) commandCount(conn redcon.Conn) {
	conn.WriteInt(len(h.commandDescriptions))
}

// COMMAND LIST
// Return array of all commands available on server.
func (h *Handler) commandList(conn redcon.Conn) {
	conn.WriteArray(len(h.commandDescriptions))
	for name := range h.commandDescriptions {
		conn.WriteString(name)
	}
}

// COMMAND INFO [<command> ...]
// Show information about given commands (if slice is nil, show all commands).
func (h *Handler) commandInfo(conn redcon.Conn, commands [][]byte) {
	if commands == nil || len(commands) == 0 {
		conn.WriteArray(len(h.commandDescriptions))
		for name, info := range h.commandDescriptions {
			writeCommandInfo(conn, name, info)
		}
		return
	}

	conn.WriteArray(len(commands))
	for _, argB := range commands {
		arg := strings.ToLower(string(argB))
		info, ok := h.commandDescriptions[arg]
		if !ok {
			conn.WriteError(fmt.Sprintf("ERR command not found '%s'", arg))
			continue
		}
		writeCommandInfo(conn, arg, info)
	}
}

// COMMAND DOCS [<command> ...]
// Show usage documentation about given commands (if slice is nil, show all commands).
func (h *Handler) commandDocs(conn redcon.Conn, commands [][]byte) {
	if commands == nil || len(commands) == 0 {
		conn.WriteArray(len(h.commandDescriptions))
		for name, info := range h.commandDescriptions {
			writeCommandDocs(conn, name, info)
		}
		return
	}

	conn.WriteArray(len(commands))
	for _, argB := range commands {
		arg := strings.ToLower(string(argB))
		info, ok := h.commandDescriptions[arg]
		if !ok {
			conn.WriteError(fmt.Sprintf("ERR command not found '%s'", arg))
			continue
		}
		writeCommandDocs(conn, arg, info)
	}
}

func writeCommandInfo(conn redcon.Conn, name string, info commandInfo) {
	const commandInfoEntries = 8

	conn.WriteArray(commandInfoEntries)

	conn.WriteString(name)         // 1
	conn.WriteInt(info.arity)      // 2
	conn.WriteAny(info.flags)      // 3
	conn.WriteInt(info.firstKey)   // 4
	conn.WriteInt(info.lastKey)    // 5
	conn.WriteInt(info.stepKey)    // 6
	conn.WriteAny(info.categories) // 7
	conn.WriteAny(info.tips)       // 8
}

func writeCommandDocs(conn redcon.Conn, name string, info commandInfo) {
	const commandDocsEntries = 2

	conn.WriteArray(commandDocsEntries)
	conn.WriteString(name)   // 1
	conn.WriteAny(info.tips) // 2
}
