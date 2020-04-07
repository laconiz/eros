package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/network/websocket"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/example/service"
	"net/url"
	"sync"
)

type Service struct {
	oceanus oceanus.Oceanus
	users   map[session.ID]*UserProxy
	service *websocket.Acceptor
	mutex   sync.RWMutex
}

func (s *Service) Init() {
	s.service.Run()
}

func (s *Service) Destroy() {
	s.service.Stop()
}

func (s *Service) onConnected(event *network.Event) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	logger := logger.Data(event.Ses.ID())

	data, ok := event.Ses.Load("ws")
	if !ok {
		logger.Error("can not load request from session")
		return
	}

	uri, ok := data.(*url.URL)
	if !ok {
		logger.Error("invalid request")
		return
	}

	query := uri.RawQuery

	id, err := s.oceanus.Create(service.User)
	if err != nil {
		logger.Err(err)
	}
}

func (s *Service) onDisconnected(event *network.Event) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.users[event.Ses.ID()]
	if !ok {
		return
	}

	s.oceanus.Destroy()
}

const module = "proxy"

var logger = logisor.Module(module)
