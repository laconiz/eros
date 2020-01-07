package oceanus

import (
	"fmt"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"time"
)

// 新建一个连接
func (p *Process) NewConnector(node *Node) network.Connector {

	invoker := network.NewStdInvoker()

	invoker.Register(network.Connected{}, func(event *network.Event) {

		p.mutex.Lock()
		defer p.mutex.Unlock()

		if burl, ok := p.burls[node.ID]; ok {
			burl.session = event.Session
		}
	})

	invoker.Register(network.Disconnected{}, func(event *network.Event) {

		p.mutex.Lock()
		defer p.mutex.Unlock()

		if burl, ok := p.burls[node.ID]; ok {
			burl.session = nil
		}
	})

	conf := tcp.ConnectorConfig{
		Name:      fmt.Sprintf("oceanus.connector.%s", node.Addr),
		Addr:      node.Addr,
		Reconnect: true,
		Session: tcp.SessionConfig{
			ReadTimeout:  time.Second * 6,
			WriteTimeout: time.Second * 6,
			LogLevel:     0,
			QueueLen:     0,
			Invoker:      nil,
			EncoderMaker: nil,
		},
	}
}

//
func (p *Process) NewAcceptor(node *Node) network.Acceptor {

	invoker := p.NewCommonInvoker()

	invoker.Register(network.Connected{})

}

func (p *Process) NewCommonInvoker() *network.StdInvoker {

	invoker := network.NewStdInvoker()

	invoker.Register(Message{}, func(event *network.Event) {
		p.Push(event.Msg.(*Message))
	})

	invoker.Register(NodeJoinMsg{}, func(event *network.Event) {
		p.OnNodeJoin(event.Msg.(*NodeJoinMsg).Nodes)
	})

	invoker.Register(NodeQuitMsg{}, func(event *network.Event) {
		p.OnNodeQuit(event.Msg.(*NodeQuitMsg).Node)
	})

	invoker.Register(ChannelJoinMsg{}, func(event *network.Event) {
		p.OnRouteJoin(event.Msg.(*ChannelJoinMsg).Channels)
	})

	invoker.Register(ChannelQuitMsg{}, func(event *network.Event) {
		p.OnRouteQuit(event.Msg.(*ChannelQuitMsg).Channels)
	})

	return invoker
}
