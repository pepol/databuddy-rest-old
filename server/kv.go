package server

import (
	"fmt"

	"github.com/pepol/databuddy/internal/log"
	"github.com/tidwall/redcon"
)

// This file contains implementation of the "kv" commands.

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

	h.mutex.RLock()
	val, err := h.db.Get(key)
	h.mutex.RUnlock()

	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR getting item '%s': %v", key, err))
		return
	}

	conn.WriteBulk(val)
}

// SET <key> <value>
// Set key to contain value.
//noling:nolintlint
//nolint:ifshort // Error shouldn't be checked and responded to inside mutex lock.
func (h *Handler) set(conn redcon.Conn, cmd redcon.Command) {
	const setArgsCount = 3

	if len(cmd.Args) != setArgsCount {
		wrongArgs(conn, string(cmd.Args[0]))
		return
	}

	key := string(cmd.Args[1])
	val := cmd.Args[2]

	// TODO: Add more argument checking.

	h.mutex.Lock()
	err := h.db.Set(key, val)
	h.mutex.Unlock()

	if err != nil {
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

	deleted := 0

	for i := 1; i < len(cmd.Args); i++ {
		key := string(cmd.Args[i])

		// TODO: Add more argument checking.

		h.mutex.Lock()
		err := h.db.Delete(key)
		h.mutex.Unlock()

		if err != nil {
			log.Error(fmt.Sprintf("deleting key '%s'", key), err)
			continue
		}
		deleted++
	}

	conn.WriteInt(deleted)
}
