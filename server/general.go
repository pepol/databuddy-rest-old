package server

import (
	"fmt"
	"os"

	"github.com/pepol/databuddy/internal/log"
	"github.com/tidwall/redcon"
)

// This file contains implementation of the "general" commands.

func (h *Handler) info(conn redcon.Conn, _cmd redcon.Command) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	conn.WriteBulkString(fmt.Sprintf(
		"DataBuddy %s %s (%s) client: %s\r\n",
		version,
		addr,
		hostname,
		conn.RemoteAddr(),
	))
}

func (h *Handler) ping(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("PONG")
}

func (h *Handler) quit(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("BYE")
	if err := conn.Close(); err != nil {
		log.Error("closing connection", err)
	}
}
