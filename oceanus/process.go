package oceanus

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/websocket"
	"net"
	"sync"
)

var (
	threads = map[string]*Thread{}

	logger = log.Std("oceanus")

	acceptor network.Acceptor

	connections = map[string]network.Connector{}
)

type Config struct {
	Addr string
}

func runAcceptor(conf Config) {

	acceptor = websocket.NewAcceptor(websocket.AcceptorConfig{
		Name: "oceanus.lookup",
		Addr: conf.Addr,
	})

	acceptor.Start()
}

func Run(conf Config) error {

	runAcceptor(conf)

	return nil
}

type Peer struct {
	channels    map[string]*Channel
	clients map[string]*net.Conn
	mutex       sync.Mutex
}

func (p *Peer) Run() error {

	listener, err := net.Listen("tcp", "192.168.10.106:4369")
	if err != nil {
		return err
	}

	for {

		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		p.handlerConn(conn)
	}
}

func (p *Peer) handlerConn(conn net.Conn) {

	p.mutex.Lock()
	p.connections[]

}
