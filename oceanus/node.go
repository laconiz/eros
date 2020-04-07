package oceanus

type NodeID string

type NodeInfo struct {
	ID      NodeID
	Routers Routers
}

// ---------------------------------------------------------------------------------------------------------------------

type Node interface {
	Info() *NodeInfo
	Mesh() Mesh
	Mail(*Mail) error
}

// ---------------------------------------------------------------------------------------------------------------------

type localNode struct {
	info *NodeInfo
}

func (node *localNode) Info() *NodeInfo {
	return node.info
}

// ---------------------------------------------------------------------------------------------------------------------

type remoteNode struct {
	info *NodeInfo
}

func (node *remoteNode) Info() *NodeInfo {
	return node.info
}
