package server

import (
	"fmt"

	"github.com/pepol/databuddy/internal/log"
	"github.com/tidwall/redcon"
)

// This file contains implementation of the "general" commands.

// NODE
// Returns node information.
func (h *Handler) nodeInfo(conn redcon.Conn, _cmd redcon.Command) {
	conn.WriteBulkString(fmt.Sprintf(
		"DataBuddy %s %s (%s) client: %s\r\n",
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
