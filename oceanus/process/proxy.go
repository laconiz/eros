package process

import (
	"errors"
	"fmt"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/channel"
	"math/rand"
)

var ErrNotFound = errors.New("can not find channel")

// 投递消息
// 注意：message的receiver有重复时只会投递一次
func (p *Process) post(message *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// 记录已经被投递过消息的节点
	sn := map[string]bool{}

	for _, receiver := range message.Receivers {

		// 已投递过的消息节点
		if sn[receiver.Node.UUID] {
			continue
		}

		// 投递消息
		if n, ok := p.nodes[receiver.Node.UUID]; ok {
			n.Send(message)
		}

		// 记录消息节点
		sn[receiver.Node.UUID] = true
	}

	return nil
}

func (p *Process) direct(typo, key string, msg *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	group, ok := p.groups[typo]
	if !ok {
		return fmt.Errorf("%w: type[%s]", ErrNotFound, typo)
	}

	channel, ok := group[key]
	if !ok {
		return fmt.Errorf("%w: channel[%s.%s]", ErrNotFound, typo, key)
	}

	return channel.Push(msg)
}

func (p *Process) multicast(typo string, keys []string, message *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	group, ok := p.groups[typo]
	if !ok || len(group) == 0 {
		return fmt.Errorf("%w: type[%s]", ErrNotFound, typo)
	}

	for _, key := range keys {

		channel, ok := group[key]
		if !ok {
			continue
		}

		// 本地通道直接投递
		if channel.Local() {
			msg := message.Copy()
			msg.Receiver = channel.Info()
			channel.Push(msg)
			continue
		}

	}

	return nil
}

func (p *Process) broadcast(typo string, message *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	group, ok := p.groups[typo]
	if !ok || len(group) == 0 {
		return fmt.Errorf("can not find type: %s", typo)
	}

	for _, channel := range group {
		msg := message.Copy()
		msg.Receiver = channel.Info()
		channel.Push(msg)
	}

	return nil
}

func (p *Process) balance(typo string, message *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	group, ok := p.groups[typo]
	if !ok || len(group) == 0 {
		return fmt.Errorf("%w: type[%s]", ErrNotFound, typo)
	}

	var channels []*channel.Channel

	index := 0
	slice := make([]oceanus.Channel, len(group))
	for _, channel := range group {
		slice[index] = channel
		index++
	}

	channel := slice[rand.Intn(len(slice))]
	msg := message.Copy()
	msg.Receiver = channel.Info()

	return channel.Push(msg)
}
