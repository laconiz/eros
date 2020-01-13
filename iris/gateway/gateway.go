package gateway

import (
	"github.com/laconiz/eros/iris/proto"
	"github.com/laconiz/eros/network"
	"sync"
)

func newGateway() *Gateway {
	return &Gateway{sessions: map[proto.UserID]network.Session{}}
}

type Gateway struct {
	sessions map[proto.UserID]network.Session
	mutex    sync.RWMutex
}

func (g *Gateway) SendUserMessage(userID proto.UserID, msg interface{}) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	if session, ok := g.sessions[userID]; ok {
		session.Send(msg)
	}
}

func (g *Gateway) SendUserStream(id proto.UserID, stream []byte) {
	g.mutex.Lock()
	defer g.mutex.RUnlock()
	if session, ok := g.sessions[id]; ok {
		session.SendStream(stream)
	}
}

func (g *Gateway) Store(id proto.UserID, session network.Session) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.sessions[id] = session
}

func (g *Gateway) Delete(id proto.UserID) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	delete(g.sessions, id)
}
