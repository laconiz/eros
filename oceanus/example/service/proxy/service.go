package proxy

import (
	"errors"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/network/websocket"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/example/model"
	"github.com/laconiz/eros/oceanus/example/service"
	"net/url"
	"sync"
)

type Service struct {
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
		logger.Error("can not load url from session")
		return
	}

	uri, ok := data.(*url.URL)
	if !ok {
		logger.Data(data).Error("invalid url data type on session")
		return
	}

	query := uri.RawQuery

	id, err := oceanus.Create(service.User)
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

// ---------------------------------------------------------------------------------------------------------------------

// 将websocket
func (s *Service) verifyConnection(ses session.Session) (model.UserID, error) {

	data, ok := ses.Load(websocket.KeyURL)
	if !ok {
		return 0, errors.New("can not load url from session")
	}

	url, ok := data.(*url.URL)
	if !ok {
		return 0, errors.New()
	}

}

// ---------------------------------------------------------------------------------------------------------------------

const module = "proxy"

var logger = logisor.Module(module)
