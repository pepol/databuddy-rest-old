package server

import (
	"fmt"

	"github.com/pepol/databuddy/internal/log"
	"github.com/tidwall/redcon"
)

// This file contains implementation of the "general" commands.

// INFO
// Returns node information.
func (h *Handler) info(conn redcon.Conn, _cmd redcon.Command) {
	conn.WriteBulkString(fmt.Sprintf(
		"DataBuddy %s %s (%s) client: %s",
		h.version,
		h.addr,
		h.hostname,
		conn.RemoteAddr(),
	))
}

// PING
// Responds with pong.
func (h *Handler) ping(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("PONG")
}

// QUIT
// Closes connection.
func (h *Handler) quit(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("BYE")
	if err := conn.Close(); err != nil {
		log.Error("closing connection", err)
	}
}

func registerGeneral(handler *Handler) {
	handler.Register("info", handler.info, 1, []string{"server"}, -1, -1, 0, nil, []string{"INFO", "return information about current node"})
	handler.Register("ping", handler.ping, 1, []string{"general"}, -1, -1, 0, nil, []string{"PING", "respond with 'PONG'"})
	handler.Register("quit", handler.quit, 1, []string{"server"}, -1, -1, 0, nil, []string{"QUIT", "close the connection"})
}
