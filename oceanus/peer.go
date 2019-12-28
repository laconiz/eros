package oceanus

import "github.com/laconiz/eros/network"

type Peer struct {
	connections map[string]network.Connector
}

func (p *Peer) Send(event *Event) error {

}

func (p *Peer)
