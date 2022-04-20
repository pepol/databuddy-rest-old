package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/serf/serf"
	"github.com/tidwall/redcon"
)

// This file contains implementation of the "cluster management" commands.

const clusterArgCount = 2

// CLUSTER
// Basic handler for cluster command container.
func (h *Handler) cluster(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != clusterArgCount {
		wrongArgs(conn, "CLUSTER")
		return
	}

	subcommand := strings.ToLower(string(cmd.Args[1]))

	switch subcommand {
	case "count": //nolint:goconst // "count" happens to be same for multiple commands, however it is independent for each.
		h.clusterCount(conn)
	case "peers":
		h.clusterPeers(conn)
	default:
		conn.WriteError(fmt.Sprintf("ERR unknown command '%s %s'", string(cmd.Args[0]), subcommand))
	}
}

// CLUSTER COUNT
// Return count of all cluster members.
func (h *Handler) clusterCount(conn redcon.Conn) {
	conn.WriteInt(h.serf.NumNodes())
}

// CLUSTER PEERS
// Return list of all cluster members.
func (h *Handler) clusterPeers(conn redcon.Conn) {
	members := h.serf.Members()
	membersCount := len(members)

	conn.WriteArray(membersCount)
	for _, member := range members {
		writePeerInfo(conn, member)
	}
}

//nolint:gomnd // Magic numbers in command registration calls are not magic.
func registerCluster(handler *Handler) {
	handler.Register("cluster", handler.cluster, 1, []string{"cluster"}, 1, 1, 0, nil, []string{"CLUSTER", "container for cluster commands"})
	handler.RegisterChild("cluster count", 2, []string{"cluster"}, -1, -1, 0, nil, []string{"CLUSTER COUNT", "return count of all known members of cluster"})
	handler.RegisterChild("cluster peers", 2, []string{"cluster"}, -1, -1, 0, nil, []string{"CLUSTER PEERS", "return list of all known members of cluster"})
}

func writePeerInfo(conn redcon.Conn, member serf.Member) {
	const peerInfoEntries = 5
	const peerVersionEntries = 3 // Min, Cur, Max

	versions := make([]uint64, peerVersionEntries)
	versions[0] = uint64(member.ProtocolMin)
	versions[1] = uint64(member.ProtocolCur)
	versions[2] = uint64(member.ProtocolMax)

	conn.WriteArray(peerInfoEntries)

	conn.WriteString(member.Name)                                                            // 1 - node name
	conn.WriteString(net.JoinHostPort(member.Addr.String(), fmt.Sprintf("%d", member.Port))) // 2 - connection string
	conn.WriteAny(member.Tags)                                                               // 3 - peer tags
	conn.WriteString(member.Status.String())                                                 // 4 - peer status
	conn.WriteAny(versions)                                                                  // 5 - versions
}
