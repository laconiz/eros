// 本地节点管理器

package local

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/oceanus/abstract"
	"github.com/laconiz/eros/oceanus/proto"
)

func New(info *proto.Mesh, proc abstract.Process) *Mesh {
	return &Mesh{
		info:    info,
		state:   &proto.State{},
		nodes:   map[proto.NodeID]*Node{},
		types:   map[proto.NodeType]int32{},
		logger:  proc.Logger(),
		process: proc,
	}
}

type Mesh struct {
	info    *proto.Mesh              // 网格信息
	state   *proto.State             // 网格状态
	nodes   map[proto.NodeID]*Node   // 节点列表
	types   map[proto.NodeType]int32 // 节点类型统计
	logger  logis.Logger             // 日志接口
	process abstract.Process
}

func (mesh *Mesh) Info() *proto.Mesh {
	return mesh.info
}

func (mesh *Mesh) State() (*proto.State, bool) {
	return mesh.state, true
}

func (mesh *Mesh) Nodes() []*proto.Node {

	var nodes []*proto.Node
	for _, node := range mesh.nodes {
		nodes = append(nodes, node.Info())
	}

	return nodes
}

func (mesh *Mesh) Mail(mail *proto.Mail) error {

	for _, ni := range mail.To {

		node, ok := mesh.nodes[ni.ID]
		if !ok {
			continue
		}

		node.Mail(mail)
	}

	return nil
}

func (mesh *Mesh) Destroy() {

	for id, node := range mesh.nodes {
		node.Destroy()
	}

	mesh.nodes = map[proto.NodeID]*Node{}
	mesh.types = map[proto.NodeType]int32{}
}

//
func (mesh *Mesh) Create(ni proto.Node, handler interface{}) *Node {

}

//
func (mesh *Mesh) Delete(id proto.NodeID) {

}
