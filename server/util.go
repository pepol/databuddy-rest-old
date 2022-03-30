package server

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func wrongArgs(conn redcon.Conn, command string) {
	conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", command))
}
