package oceanus

type MeshID string

type MeshInfo struct {
	ID   MeshID
	Addr string
}

// ---------------------------------------------------------------------------------------------------------------------

type Mesh interface {
	Info() *MeshInfo
	Mail(*Mail) error
}

// ---------------------------------------------------------------------------------------------------------------------

type localMesh struct {
	info  *MeshInfo
	nodes map[NodeID]*localNode
}

func (mesh *localMesh) Info() *MeshInfo {
	return mesh.info
}

func (mesh *localMesh) Nodes() []*NodeInfo {
	var nodes []*NodeInfo
	for _, node := range mesh.nodes {
		nodes = append(nodes, node.info)
	}
	return nodes
}

// ---------------------------------------------------------------------------------------------------------------------

type remoteMesh struct {
	info  *MeshInfo
	nodes map[NodeID]*remoteNode
}

func (mesh *remoteMesh) Info() *MeshInfo {
	return mesh.info
}

func (mesh *remoteMesh) Nodes() []*NodeInfo {
	var nodes []*NodeInfo
	for _, node := range mesh.nodes {
		nodes = append(nodes, node.info)
	}
	return nodes
}
