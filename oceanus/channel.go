package oceanus

import (
	"github.com/laconiz/eros/queue"
)

type ChannelID string

type ChannelType string

type ChannelKey string

type Channel interface {
	Info() *ChannelInfo
	Push(*Message) error
}

type ChannelInfo struct {
	ID   ChannelID
	Type ChannelType
	Key  ChannelKey
}

func (c *ChannelInfo) Info() *ChannelInfo {
	return c
}

type remoteChannel struct {
	*ChannelInfo
	node Node
}

func (c *remoteChannel) Push(message *Message) error {
	return c.node.Push(message)
}

func newRemoteChannel(info *ChannelInfo, node Node) Channel {
	return &remoteChannel{ChannelInfo: info, node: node}
}

type localChannel struct {
	*ChannelInfo
	node  Node
	queue *queue.Queue
}

func (c *localChannel) Push(message *Message) error {
	return c.queue.Add(message)
}

func newLocalChannel(info *ChannelInfo, node Node) Channel {
	return &localChannel{ChannelInfo: info, node: node}
}
