package v1alpha3

// ClusterLeader is a connection string for the current Raft cluster leader.
type ClusterLeader string

// ClusterPeers is an array of connection strings for all peers joined in cluster.
type ClusterPeers []string

// ClusterService is a service for management of clusters.
type ClusterService interface {
	Local() *ClusterInfo
	Get(string) (*ClusterInfo, error)
}

// ClusterInfo contains all the information about given cluster.
type ClusterInfo struct {
	Name        string
	Leader      string
	Peers       []string
	Annotations map[string]string
	Labels      []string
}
