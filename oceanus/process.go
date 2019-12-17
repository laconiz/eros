package oceanus

import (
	"github.com/laconiz/eros/network"
	"sync"
)

var channels = map[string]*Channel{}

var channelByGroup = map[string]map[string]Channel{}

var mutex sync.RWMutex

func Direct(typo, key string, raw []byte) {

	mutex.RLock()
	defer mutex.RUnlock()

	if group, ok := channelByGroup[typo]; ok {
		if channel, ok := group[key]; ok {
			channel.
		}
	}
}

func Group(typo string, keys []string, raw []byte) {

}

func Broadcast(typo string, raw []byte) {

	messages := map[string]*Message{}

	for key, channel := range channelByGroup[typo] {
		receivers := messages[channel.Info().Peer].Receivers
		receivers = append(receivers, key)
		messages[channel.Info().Peer].Receivers = receivers
	}

}

type Config struct {
	Addr string
}

type Peer struct {
	channels   map[string]*Channel
	states     map[string]*ChannelStates
	acceptor   network.Acceptor
	connectors map[string]network.Connector

	mutex sync.Mutex
}

func (p *Peer) Run() {

}

func (p *Peer) NewChannel() {

	delivery := make()

}

func init() {

}
