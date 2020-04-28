// 远程网格
package remote

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/abstract"
	"github.com/laconiz/eros/oceanus/proto"
	"hash/fnv"
)

// 创建远程网格
func New(info *proto.Mesh, proc abstract.Process) *Mesh {

	mesh := &Mesh{
		info:  info,
		state: &proto.State{},
		nodes: map[proto.NodeID]*Node{},
		types: map[proto.NodeType]int32{},
	}

	// 计算ADDR HASH值
	hash := fnv.New32()
	hash.Write([]byte(info.Addr))
	rp := hash.Sum32()
	hash.Reset()
	hash.Write([]byte(proc.Local().Info().Addr))
	lp := hash.Sum32()

	// 创建客户端连接
	if lp >= rp && (lp-rp)%2 == 0 ||
		rp > lp && (rp-lp)%2 != 0 {

		proc.Logger().Data(info.Addr).Info("connect to")
		mesh.connector = proc.NewConnector(info.Addr)
		go mesh.connector.Run()
	}

	return mesh
}

// 远程网格
type Mesh struct {
	info      *proto.Mesh              // 网格信息
	state     *proto.State             // 网格状态
	nodes     map[proto.NodeID]*Node   // 节点列表
	types     map[proto.NodeType]int32 // 节点类型统计
	session   session.Session          // 网络连接
	connector network.Connector        // 连接器
	router    abstract.Router          // 路由器
}

// 网格信息
func (mesh *Mesh) Info() *proto.Mesh {
	return mesh.info
}

// 网格状态
func (mesh *Mesh) State() (*proto.State, bool) {
	return mesh.state, mesh.session != nil
}

// 节点列表
func (mesh *Mesh) Nodes() []*proto.Node {

	var nodes []*proto.Node
	for _, node := range mesh.nodes {
		nodes = append(nodes, node.Info())
	}

	return nodes
}

// 发送邮件
func (mesh *Mesh) Mail(mail *proto.Mail) error {

	if mesh.session == nil {
		return oceanus.ErrDisconnected
	}

	return mesh.session.Send(mail)
}

// 发送消息
func (mesh *Mesh) Send(msg interface{}) error {

	if mesh.session == nil {
		return oceanus.ErrDisconnected
	}

	return mesh.session.Send(msg)
}

// 更新网格状态
func (mesh *Mesh) UpdateState(state *proto.State) {
	mesh.state = state
	mesh.expired()
}

// 更新网格连接
func (mesh *Mesh) UpdateSession(session session.Session) {
	mesh.session = session
	mesh.expired()
}

// 设置路由器过期
func (mesh *Mesh) expired() {

	for typo, count := range mesh.types {

		if count == 0 {
			continue
		}

		mesh.router.Expired(typo)
	}
}

// 销毁网格
func (mesh *Mesh) Destroy() {

	for _, node := range mesh.nodes {
		node.Destroy()
	}

	// 删除连接
	if mesh.connector != nil {
		mesh.connector.Stop()
		mesh.connector = nil
	}

	mesh.nodes = map[proto.NodeID]*Node{}
	mesh.types = map[proto.NodeType]int32{}
}

// 插入节点
func (mesh *Mesh) Insert(nodes []*proto.Node) {

	mesh.Remove(nodes)

	for _, info := range nodes {
		node := newNode(info, mesh)
		mesh.router.Insert(node)
		mesh.types[info.Type]++
	}
}

// 移除节点
func (mesh *Mesh) Remove(nodes []*proto.Node) {

	for _, info := range nodes {

		_, ok := mesh.nodes[info.ID]
		if !ok {
			continue
		}

		mesh.router.Remove(info.ID)
		mesh.types[info.Type]--
		delete(mesh.nodes, info.ID)
	}
}
