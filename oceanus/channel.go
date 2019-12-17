package oceanus

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/queue"
)

type Channel interface {
	Info() *ChannelInfo
	State() *ChannelState
}

type LocalChannel struct {
	Info  ChannelInfo
	State ChannelState
	Queue *queue.Queue
}

func (c *LocalChannel) Send(message Message) {

}

type ChannelInfo struct {
	UUID string
	Type string
	Name string
	Peer string
}

type ChannelState struct {
	UUID string
}

type Message struct {
	Receivers []string
	Sender    string
	Body      []byte
}

type Channels struct {
	Channels []*ChannelInfo
}

type ChannelStates struct {
	States []*ChannelState
}

func init() {
	network.RegisterMeta(Message{}, network.JsonCodec)
	network.RegisterMeta(Channels{}, network.JsonCodec)
	network.RegisterMeta(ChannelStates{}, network.JsonCodec)
}
