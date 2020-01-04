package oceanus

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"sync"
	"time"
)

type NodeID string

type NodeState struct {
	Version uint32
}

type Node interface {
	Info() *NodeInfo
	Push(*Message) error
}

type NodeInfo struct {
	ID    NodeID
	Addr  string
	State NodeState
}

func (n *NodeInfo) Info() *NodeInfo {
	return n
}

type remoteNode struct {
	NodeInfo
	channels []Channel
	conn     network.Connector
}

func (n *remoteNode) Push(message *Message) error {
	return n.conn.Send(message)
}

func NewRemoteNode(info *NodeInfo) Node {
	return &remoteNode{NodeInfo: *info}
}

type thread struct {
	*NodeInfo
	channels map[ChannelID]Channel
	acceptor network.Acceptor
	remotes  map[NodeID]*remoteNode
	mutex    sync.RWMutex
}

func (n *thread) Push(message *Message) error {

	n.mutex.RLock()
	defer n.mutex.RUnlock()

	for _, id := range message.Receivers {

		if channel, ok := n.channels[id]; ok {
			channel.Push(message)
		}
	}

	return nil
}

func (n *thread) joinNode(info NodeInfo) {

	n.mutex.Lock()
	defer n.mutex.Unlock()

	node, ok := n.remotes[info.ID]
	if ok && !node.conn.Connected() {
		node.conn.Run()
	}

}

func (n *thread) invoker() network.Invoker {

	invoker := network.NewStdInvoker()

	invoker.Register(&Message{}, func(event *network.Event) {
		n.Push(event.Msg.(*Message))
	})

	invoker.Register(NodeJoinMsg{}, func(event *network.Event) {

		msg := event.Msg.(*NodeJoinMsg)

		n.mutex.Lock()
		defer n.mutex.Unlock()

		node, ok := n.remotes[msg.Node.ID]
		if ok {
			if node
		}

		node = NewRemoteNode(&msg.Node)
		n.remotes[msg.Node.ID] = node

		n.mutex.RLock()
		defer n.mutex.RUnlock()

		msg := &NodeListMsg{}

		func

		for _, node := range n.remotes {

		}
	})
}

func newLocalNode(info *NodeInfo) Node {

	node := &thread{
		NodeInfo: info,
		channels: map[ChannelID]Channel{},
	}

	conf := tcp.AcceptorConfig{
		Name:    "oceanus.node",
		Addr:    info.Addr,
		Session: nodeSessionConfig,
	}
	conf.Session.Invoker = invoker
	node.acceptor = tcp.NewAcceptor(conf)

	return node
}

var nodeSessionConfig = tcp.SessionConfig{
	ReadTimeout:  time.Second * 6,
	WriteTimeout: time.Second * 6,
	LogLevel:     log.Warn,
	QueueLen:     64,
}
