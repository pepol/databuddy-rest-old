package server

import (
	"fmt"
	"strings"

	"github.com/pepol/databuddy/internal/context"
	"github.com/pepol/databuddy/internal/log"
	"github.com/tidwall/redcon"
)

// This file contains implementation of the "database management" commands.

// BUCKET
// Return name of currently used bucket.
func (h *Handler) bucket(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) == 1 {
		ctx, ok := conn.Context().(*context.Context)
		if !ok {
			conn.WriteError("ERR context not set on connection")
			if err := conn.Close(); err != nil {
				log.Error("closing connection", err)
			}
			return
		}

		conn.WriteString(ctx.Bucket.Name)
		return
	}

	subcommand := strings.ToLower(string(cmd.Args[1]))

	switch subcommand {
	case "count":
		h.bucketCount(conn)
	case "list":
		h.bucketList(conn, cmd.Args[2:])
	case "create":
		h.bucketCreate(conn, cmd.Args[2:])
	case "use":
		h.bucketUse(conn, cmd.Args[2:])
	case "drop":
		h.bucketDrop(conn, cmd.Args[2:])
	default:
		conn.WriteError(fmt.Sprintf("ERR unknown command '%s %s'", string(cmd.Args[0]), subcommand))
	}
}

// BUCKET COUNT
// Return count of all buckets available on server.
func (h *Handler) bucketCount(conn redcon.Conn) {
	conn.WriteInt(h.db.Count())
}

// BUCKET LIST
// Return array of all buckets available on server.
func (h *Handler) bucketList(conn redcon.Conn, args [][]byte) {
	var prefix string

	switch len(args) {
	case 0:
		prefix = ""
	case 1:
		prefix = strings.ToLower(string(args[0]))
	default:
		wrongArgs(conn, "BUCKET LIST")
		return
	}

	conn.WriteAny(h.db.List(prefix))
}

// BUCKET CREATE <bucket>
// Create bucket with given name.
func (h *Handler) bucketCreate(conn redcon.Conn, args [][]byte) {
	if len(args) != 1 {
		wrongArgs(conn, "BUCKET CREATE")
		return
	}

	name := string(args[0])

	if err := h.db.Create(name); err != nil {
		conn.WriteError(fmt.Sprintf("ERR creating bucket '%s': %v", name, err))
		return
	}

	conn.WriteString("OK")
}

// BUCKET USE <bucket>
// Set bucket for further queries.
func (h *Handler) bucketUse(conn redcon.Conn, args [][]byte) {
	if len(args) != 1 {
		wrongArgs(conn, "BUCKET USE")
		return
	}

	name := string(args[0])

	bucket, err := h.db.Get(name)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR opening bucket '%s': %v", name, err))
		return
	}

	ctx, ok := conn.Context().(*context.Context)
	if !ok {
		ctx = &context.Context{
			Bucket: bucket,
		}
		conn.SetContext(ctx)
		conn.WriteString("OK context set (not set previously)")
		return
	}

	if name == ctx.Bucket.Name {
		conn.WriteString("OK bucket already used")
		return
	}

	ctx.Bucket = bucket

	conn.WriteString("OK")
}

// BUCKET DROP <bucket> [<bucket> ...]
// Remove given buckets, including all data.
func (h *Handler) bucketDrop(conn redcon.Conn, args [][]byte) {
	if len(args) == 0 {
		wrongArgs(conn, "BUCKET DROP")
		return
	}

	ctx, ok := conn.Context().(*context.Context)
	if !ok {
		conn.WriteError("ERR context not set on connection")
		return
	}

	dropped := 0

	for _, arg := range args {
		name := string(arg)

		// TODO: Improve check with all currently connected users.
		if name == ctx.Bucket.Name {
			log.Warn(fmt.Sprintf("not dropping bucket '%s' as its used by current user", name))
			continue
		}

		if err := h.db.Drop(name); err != nil {
			log.Error(fmt.Sprintf("dropping bucket '%s'", name), err)
			continue
		}
		dropped++
	}

	conn.WriteInt(dropped)
}
