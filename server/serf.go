package server

import (
	"log"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
)

func getSerfConfig(host string, port int, id string, logger *log.Logger, eventCh chan<- serf.Event) *serf.Config {
	memberlistConfig := memberlist.DefaultLANConfig()
	memberlistConfig.BindAddr = host
	memberlistConfig.BindPort = port
	memberlistConfig.Logger = logger
	memberlistConfig.ProtocolVersion = 5

	serfConfig := serf.DefaultConfig()
	serfConfig.Tags = map[string]string{
		"role": "kv",
	}
	serfConfig.NodeName = id
	serfConfig.EventCh = eventCh
	serfConfig.MemberlistConfig = memberlistConfig
	serfConfig.Logger = logger
	serfConfig.ProtocolVersion = 5

	return serfConfig
}
