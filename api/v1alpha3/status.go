package v1alpha3

import (
	"github.com/gofiber/fiber/v2"
)

// NodeService provides information about current node.
type NodeService interface {
	Status() (*NodeInfo, error)
}

// NodeInfo contains all the information about current node.
type NodeInfo struct {
	Hostname    string
	Version     string
	CPU         []NodeCPUInfo
	Memory      NodeMemoryInfo
	Disk        []NodeDiskInfo
	Cluster     string
	Labels      []string
	Annotations map[string]string
}

// NodeCPUInfo contains CPU information for node.
type NodeCPUInfo struct {
	CPUIndex  int32
	VendorID  string
	Family    string
	Model     string
	Cores     int32
	ModelName string
	Mhz       float64
	CacheSize int32
	Flags     []string
	Microcode string
}

// NodeMemoryInfo contains memory information for node.
type NodeMemoryInfo struct {
	Total       uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
	Free        uint64
}

// NodeDiskInfo contains disk information for node.
type NodeDiskInfo struct {
	Path              string
	Fstype            string
	Total             uint64
	Free              uint64
	Used              uint64
	UsedPercent       float64
	InodesTotal       uint64
	InodesUsed        uint64
	InodesFree        uint64
	InodesUsedPercent float64
}

func (ctl *Controller) routeStatus(router fiber.Router) {
	status := router.Group("/status")

	node := status.Group("/node")
	node.Get("/", ctl.getNodeStatus)

	cluster := status.Group("/cluster")
	cluster.Get("/leader", ctl.getClusterLeader)
	cluster.Get("/peers", ctl.getClusterPeers)
}

// @Summary Get node status
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} NodeInfo
// @Router /status/node [get]
// @Description Get current node information and status.
func (ctl *Controller) getNodeStatus(c *fiber.Ctx) error {
	status, err := ctl.self.Status()
	if err != nil {
		return err
	}

	return c.JSON(status)
}

// @Summary Get Raft leader
// @Tags status
// @Accept json
// @Produce json
// @Param dc query string false "Datacenter to query, defaults to local datacenter" default(local)
// @Success 200 {object} ClusterLeader
// @Failure 404 {object} RequestError "Requested datacenter doesn't exist"
// @Router /status/cluster/leader [get]
// @Description Get cluster's Raft leader information.
func (ctl *Controller) getClusterLeader(c *fiber.Ctx) error {
	dc := c.Query("dc", "local")

	if dc == "local" || dc == ctl.cluster.Local().Name {
		return c.JSON(ctl.cluster.Local().Leader)
	}

	clusterInfo, err := ctl.cluster.Get(dc)
	if err != nil {
		return httpError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(clusterInfo.Leader)
}

// @Summary Get cluster peers
// @Tags status
// @Accept json
// @Produce json
// @Param dc query string false "Datacenter to query, defaults to local datacenter" default(local)
// @Success 200 {object} ClusterPeers
// @Failure 404 {object} RequestError "Requested datacenter doesn't exist"
// @Router /status/cluster/peers [get]
// @Description Get the peers for given cluster.
func (ctl *Controller) getClusterPeers(c *fiber.Ctx) error {
	dc := c.Query("dc", "local")

	if dc == "local" || dc == ctl.cluster.Local().Name {
		return c.JSON(ctl.cluster.Local().Peers)
	}

	clusterInfo, err := ctl.cluster.Get(dc)
	if err != nil {
		return httpError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(clusterInfo.Peers)
}
