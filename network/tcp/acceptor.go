package tcp

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"net"
	"sync"
)

type AcceptorConfig struct {
	Addr string
}

type Acceptor struct {
	conf       *AcceptorConfig
	listener   net.Listener
	sessionMgr network.SessionMgr
	logger     *log.Logger
	mutex      sync.Mutex
}

func (acc *Acceptor) Run() {

	acc.mutex.Lock()
	if acc.listener != nil {
		acc.mutex.Unlock()
		return
	}

	listener, err := net.Listen("tcp", acc.conf.Addr)
	if err != nil {
		acc.logger.Errorf("listen at %s error: %v", acc.conf.Addr, err)
		return
	}

	acc.logger.Infof("running at: %s", acc.conf.Addr)
	acc.listener = listener
	acc.mutex.Unlock()

	for {

		conn, err := listener.Accept()
		if err != nil {
			acc.logger.Errorf("accept error: %v", err)
			break
		}
	}

	acc.mutex.Lock()
	acc.listener = nil
	acc.mutex.Unlock()
}

func (acc *Acceptor) Stop() {

	acc.mutex.Lock()
	acc.mutex.Unlock()

	if acc.listener != nil {
		acc.listener.Close()
	}
}

func NewAcceptor(conf AcceptorConfig) network.Acceptor {

}
