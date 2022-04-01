package server

import (
	"fmt"

	"github.com/pepol/databuddy/internal/context"
	"github.com/pepol/databuddy/internal/log"
	"github.com/tidwall/redcon"
)

// This file contains implementation of the "database management" commands.

// CREATE <bucket>
// Create bucket with given name.
func (h *Handler) create(conn redcon.Conn, cmd redcon.Command) {
	const createArgsCount = 2

	if len(cmd.Args) != createArgsCount {
		wrongArgs(conn, string(cmd.Args[0]))
	}

	name := string(cmd.Args[1])

	// TODO: Add more argument checking.

	if err := h.db.Create(name); err != nil {
		conn.WriteError(fmt.Sprintf("ERR creating bucket '%s': %v", name, err))
		return
	}

	conn.WriteString("OK")
}

// USE <bucket>
// Set bucket for further queries.
func (h *Handler) use(conn redcon.Conn, cmd redcon.Command) {
	const useArgsCount = 2
	if len(cmd.Args) != useArgsCount {
		wrongArgs(conn, string(cmd.Args[0]))
	}

	name := string(cmd.Args[1])

	ctx, ok := conn.Context().(*context.Context)
	if !ok {
		conn.WriteError("ERR context not set on connection")
		if err := conn.Close(); err != nil {
			log.Error("closing connection", err)
		}
		return
	}

	if name == ctx.Bucket.Name {
		conn.WriteString("OK bucket already used")
		return
	}

	bucket, err := h.db.Get(name)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR opening bucket '%s': %v", name, err))
		return
	}

	ctx.Bucket = bucket

	conn.WriteString("OK")
}
