package server

import (
	"crypto/sha256"
	"fmt"
	stdlog "log"

	"github.com/pepol/databuddy/internal/log"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/tidwall/redcon"
)

func wrongArgs(conn redcon.Conn, command string) {
	conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", command))
}

func getHostID(hostname, addr string) string {
	// Generate host ID - SHA256 hash of hostname and listen address.
	h := sha256.New()
	h.Write([]byte(hostname + addr))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func setupLogging() *stdlog.Logger {
	devel := viper.GetBool("devel")
	level := viper.GetString("loglevel")

	logger := stdlog.Default()

	l, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}

	zlog.Logger = log.Logger(l, devel)

	return logger
}
