package process

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
	"time"
)

// 序列化邮件
func (proc *Process) serialize(pkg *oceanus.Pkg) (*proto.Mail, error) {

	msg, err := proc.encoder.Marshal(pkg.Message)
	if err != nil {
		return nil, err
	}

	mailID := proto.MailID(NewGlobalUUID())

	return &proto.Mail{
		ID:      mailID,
		Headers: pkg.Headers,
		Body:    msg.Raw,
	}, nil
}

// 获取单个交换机
func (proc *Process) exchanger(pkg *oceanus.Pkg) (router.Node, error) {

	if len(pkg.IDs) > 0 {
		return proc.router.ByID(pkg.IDs[0]), nil
	}

	if pkg.Type == proto.EmptyNodeType {
		return nil, oceanus.ErrInvalidSwitch
	}

	if len(pkg.Keys) > 0 {
		return proc.router.ByKey(pkg.Type, pkg.Keys[0]), nil
	}

	return proc.router.ByLoad(pkg.Type), nil
}

// 获取一组交换机
func (proc *Process) exchangers(pkg *oceanus.Pkg) ([]router.Node, error) {

	if len(pkg.IDs) > 0 {
		return proc.router.ByIDList(pkg.IDs), nil
	}

	if pkg.Type == proto.EmptyNodeType {
		return nil, oceanus.ErrInvalidSwitch
	}

	return proc.router.ByKeys(pkg.Type, pkg.Keys), nil
}

// ---------------------------------------------------------------------------------------------------------------------

// 发送邮件
func (proc *Process) Send(pkg *oceanus.Pkg) error {

	// 序列化邮件
	mail, err := proc.serialize(pkg)
	if err != nil {
		return err
	}

	proc.mutex.Lock()
	defer proc.mutex.RLock()

	// 获取交换机
	node, err := proc.exchanger(pkg)
	if err != nil {
		return err
	}
	// 未找到对应节点
	if node == nil {
		return oceanus.ErrNodeNotFound
	}

	// 设置目标节点
	mail.To = []*proto.Node{node.Info()}

	// 发送邮件
	return node.Mail(mail)
}

// 发送一组邮件
func (proc *Process) Sends(pkg *oceanus.Pkg) error {

	// 序列化邮件
	mail, err := proc.serialize(pkg)
	if err != nil {
		return err
	}

	proc.mutex.Lock()
	defer proc.mutex.RLock()

	// 获取交换机
	nodes, err := proc.exchangers(pkg)
	if err != nil {
		return err
	}

	type Hub struct {
		mesh  router.Mesh
		nodes []*proto.Node
	}

	// 设置总线
	hubs := map[proto.MeshID]*Hub{}
	for _, node := range nodes {

		mesh := node.Mesh()
		meshID := mesh.Info().ID

		hub, ok := hubs[meshID]
		if !ok {
			hub = &Hub{mesh: mesh}
			hubs[meshID] = hub
		}
		hub.nodes = append(hub.nodes, node.Info())
	}

	// 发送邮件
	for _, hub := range hubs {
		copy := *mail
		copy.To = hub.nodes
		hub.mesh.Mail(&copy)
	}

	return nil
}

// RPC调用
func (proc *Process) Call(pkg *oceanus.Pkg) (interface{}, error) {

	// 序列化邮件
	mail, err := proc.serialize(pkg)
	if err != nil {
		return nil, err
	}
	mail.Reply = proto.RpcID(NewGlobalUUID())

	proc.mutex.Lock()

	// 获取交换机
	node, err := proc.exchanger(pkg)
	if err != nil {
		proc.mutex.Unlock()
		return nil, err
	}
	// 未找到对应节点
	if node == nil {
		proc.mutex.Unlock()
		return nil, oceanus.ErrNodeNotFound
	}

	// 设置目标节点
	mail.To = []*proto.Node{node.Info()}

	// 创建接收CHANNEL
	ch := make(chan interface{})
	proc.channels[mail.Reply] = ch

	// 发送邮件
	if err = node.Mail(mail); err != nil {
		delete(proc.channels, mail.Reply)
		proc.mutex.Unlock()
		return nil, err
	}

	proc.mutex.Unlock()

	// 等待回调消息
	var resp interface{}
	select {
	case resp = <-proc.channels[mail.Reply]:
	case <-time.After(time.Second):
		err = oceanus.ErrTimeout
	}

	// 清理接收CHANNEL
	proc.mutex.Lock()
	defer proc.mutex.Unlock()
	delete(proc.channels, mail.Reply)

	return resp, err
}
