package server

import (
	"fmt"

	"github.com/pepol/databuddy/internal/context"
	"github.com/pepol/databuddy/internal/log"
	"github.com/tidwall/redcon"
)

// This file contains implementation of the "kv" commands.

// KEYS [<prefix>]
// Get keys matching prefix (or all keys if prefix not set).
func (h *Handler) keys(conn redcon.Conn, cmd redcon.Command) {
	const keysArgsMaxCount = 2

	if len(cmd.Args) > keysArgsMaxCount {
		wrongArgs(conn, string(cmd.Args[0]))
		return
	}

	var prefix string

	if len(cmd.Args) == 1 {
		prefix = ""
	} else {
		prefix = string(cmd.Args[1])
	}

	ctx, ok := conn.Context().(*context.Context)
	if !ok {
		conn.WriteError("ERR context not set on connection")
		if err := conn.Close(); err != nil {
			log.Error("closing connection", err)
		}
		return
	}

	keys, err := ctx.Bucket.List(prefix)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR getting keys for prefix '%s': %v", prefix, err))
		return
	}

	conn.WriteAny(keys)
}

// GET <key>
// Get value at key.
func (h *Handler) get(conn redcon.Conn, cmd redcon.Command) {
	const getArgsCount = 2

	if len(cmd.Args) != getArgsCount {
		wrongArgs(conn, string(cmd.Args[0]))
		return
	}

	key := string(cmd.Args[1])

	// TODO: Add more argument checking.

	ctx, ok := conn.Context().(*context.Context)
	if !ok {
		conn.WriteError("ERR context not set on connection")
		if err := conn.Close(); err != nil {
			log.Error("closing connection", err)
		}
		return
	}

	val, err := ctx.Bucket.Get(key)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR getting item '%s': %v", key, err))
		return
	}

	conn.WriteBulk(val)
}

// SET <key> <value>
// Set key to contain value.
func (h *Handler) set(conn redcon.Conn, cmd redcon.Command) {
	const setArgsCount = 3

	if len(cmd.Args) != setArgsCount {
		wrongArgs(conn, string(cmd.Args[0]))
		return
	}

	key := string(cmd.Args[1])
	val := cmd.Args[2]

	// TODO: Add more argument checking.

	ctx, ok := conn.Context().(*context.Context)
	if !ok {
		conn.WriteError("ERR context not set on connection")
		if err := conn.Close(); err != nil {
			log.Error("closing connection", err)
		}
		return
	}

	if err := ctx.Bucket.Set(key, val); err != nil {
		conn.WriteError(fmt.Sprintf("ERR setting item '%s': %v", key, err))
		return
	}

	conn.WriteString("OK")
}

// DEL <key> [<key> ...]
// Delete key(s).
func (h *Handler) del(conn redcon.Conn, cmd redcon.Command) {
	const delArgsMinCount = 2

	if len(cmd.Args) < delArgsMinCount {
		wrongArgs(conn, string(cmd.Args[0]))
		return
	}

	ctx, ok := conn.Context().(*context.Context)
	if !ok {
		conn.WriteError("ERR context not set on connection")
		if err := conn.Close(); err != nil {
			log.Error("closing connection", err)
		}
		return
	}

	deleted := 0

	for i := 1; i < len(cmd.Args); i++ {
		key := string(cmd.Args[i])

		// TODO: Add more argument checking.

		if err := ctx.Bucket.Delete(key); err != nil {
			log.Error(fmt.Sprintf("deleting key '%s'", key), err)
			continue
		}
		deleted++
	}

	conn.WriteInt(deleted)
}
