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

	"github.com/hashicorp/serf/serf"
	"github.com/pepol/databuddy/internal/context"
	"github.com/pepol/databuddy/internal/db"
	"github.com/pepol/databuddy/internal/log"
	"github.com/spf13/viper"
	"github.com/tidwall/redcon"
)

const serfEventsBufSize = 16

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
	id        string
	accepting bool
	db        *db.Database

	serf     *serf.Serf
	eventsCh chan serf.Event

	// Replace with sorted map implementation for consistent ordering.
	commandDescriptions map[string]commandInfo

	Mux    *redcon.ServeMux
	Server *redcon.Server

	addr     string
	hostname string
	version  string
}

// NewHandler initialized the server Handler.
func NewHandler(version, addr, hostname, datadir string, s *serf.Serf, eventsCh chan serf.Event) (*Handler, error) {
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
		serf:                s,
		eventsCh:            eventsCh,
	}, nil
}

// Serve the database over network.
//nolint:funlen // Only setup boilerplate in this function, no logic.
func Serve(version string) {
	logger := setupLogging()

	port := viper.GetInt("port")
	host := viper.GetString("host")
	datadir := viper.GetString("datadir")
	join := viper.GetStringSlice("join")
	serfPort := viper.GetInt("serfport")

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	hostname, err := os.Hostname()
	if err != nil {
		log.Error("getting hostname", err)
		hostname = "localhost"
	}

	serfEvents := make(chan serf.Event, serfEventsBufSize)

	hostID := getHostID(hostname, addr)

	serfConfig := getSerfConfig(host, serfPort, hostID, logger, serfEvents)

	s, err := serf.Create(serfConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(join) > 0 {
		_, err = s.Join(join, false)
		if err != nil {
			log.Fatal(err)
		}
	}

	handler, err := NewHandler(version, addr, hostname, datadir, s, serfEvents)
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

	// Cluster commands.
	registerCluster(handler)

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

	go handler.handleSerf()

	err = server.ListenAndServe()
	if err != nil {
		log.Error("serving resp", err)
	}
	<-done
	if err := handler.serf.Shutdown(); err != nil {
		log.Error("shutting down serf", err)
	}
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

	if err := h.serf.Leave(); err != nil {
		log.Error("stopping serf", err)
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

// Handle Serf events.
func (h *Handler) handleSerf() {
	for {
		select {
		case ev := <-h.eventsCh:
			h.handleSerfEvent(ev)
		}
	}
}

func (h *Handler) handleSerfEvent(event serf.Event) {
	switch ev := event.(type) {
	case serf.MemberEvent:
		log.Info("%s: %v", ev.EventType().String(), ev.Members)
	case serf.UserEvent:
		log.Info("User: %s %v", ev.Name, ev.Payload)
	case *serf.Query:
		log.Info("Query (due at %v): %s %v", ev.Deadline(), ev.Name, ev.Payload)
		if err := ev.Respond(nil); err != nil {
			log.Error("responding to query", err)
		}
	default:
		log.Warn("unknown type: %s", ev.EventType().String())
	}
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
