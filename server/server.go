// Package server contains the main server code, including RESP handling
// and high-level database operations.
package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	accepting bool
	db        *db.Database

	// Replace with sorted map implementation for consistent ordering.
	commandDescriptions map[string]commandInfo

	Mux    *redcon.ServeMux
	Server *redcon.Server

	addr     string
	hostname string
	version  string
}

// NewHandler initialized the server Handler.
func NewHandler(version, addr, hostname, datadir string) (*Handler, error) {
	dbs, err := db.OpenDatabase(datadir)
	if err != nil {
		return nil, err
	}

	return &Handler{
		accepting:           true,
		commandDescriptions: make(map[string]commandInfo),
		db:                  dbs,
		Mux:                 redcon.NewServeMux(),
		addr:                addr,
		hostname:            hostname,
		version:             version,
	}, nil
}

// Serve the database over network.
func Serve(version string) {
	port := viper.GetInt("port")
	host := viper.GetString("host")
	datadir := viper.GetString("datadir")

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	hostname, err := os.Hostname()
	if err != nil {
		log.Error("getting hostname", err)
		hostname = "localhost"
	}

	handler, err := NewHandler(version, addr, hostname, datadir)
	if err != nil {
		log.Fatal(err)
	}

	// Meta (command-handling) commands.
	registerMeta(handler)

	// General information commands.
	registerGeneral(handler)

	// DB management commands.
	registerDatabaseManagement(handler)

	// KV commands.
	registerKV(handler)

	log.Info(fmt.Sprintf("Starting DataBuddy %s RESP server on %s", version, addr))

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	server := redcon.NewServer(
		addr,
		handler.Mux.ServeRESP,
		handler.acceptConnection,
		func(conn redcon.Conn, err error) {},
	)
	handler.Server = server

	go func() {
		<-sigs
		handler.Stop(done)
	}()

	err = server.ListenAndServe()
	if err != nil {
		log.Error("serving resp", err)
	}
	<-done
	log.Info("quitting")
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

// Stop handling connection and close database.
func (h *Handler) Stop(done chan bool) {
	h.accepting = false

	errored := false

	if err := h.Server.Close(); err != nil {
		log.Error("stopping server", err)
		errored = true
	}

	if err := h.db.Close(); err != nil {
		log.Error("closing database", err)
		errored = true
	}

	if errored {
		os.Exit(1)
	}
	done <- true
}

// Initialize connection context on connection accept.
func (h *Handler) acceptConnection(conn redcon.Conn) bool {
	if !h.accepting {
		return false
	}

	bucket, err := h.db.Get(h.db.DefaultBucket)
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

//nolint:gomnd // Magic numbers in command registration calls are not magic.
func registerMeta(handler *Handler) {
	handler.Register("command", handler.command, 1, []string{"server"}, -1, -1, 0, nil, []string{"COMMAND", "show information about all commands"})
	handler.RegisterChild("command count", 2, []string{"server"}, -1, -1, 0, nil, []string{"COMMAND COUNT", "return count of all available commands"})
	handler.RegisterChild("command list", 2, []string{"server"}, -1, -1, 0, nil, []string{"COMMAND LIST", "return list of all available commands"})
	handler.RegisterChild("command info", -2, []string{"server"}, 2, -1, 1, nil, []string{"COMMAND INFO [<command> ...]", "show system information about given command(s)"})
	handler.RegisterChild("command docs", -2, []string{"server"}, 2, -1, 1, nil, []string{"COMMAND DOCS [<command> ...]", "show documentation about given command(s)"})
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
