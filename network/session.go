package network

import (
	"fmt"
	"math"
	"sync"
)

type SessionID uint64

type Session interface {
	ID() SessionID
	Addr() string
	Send(msg interface{}) error
	SendStream(stream []byte) error
	Close()
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
}

type emptySession struct {
	values sync.Map
}

func (ses *emptySession) ID() SessionID {
	return math.MaxUint64
}

func (ses *emptySession) Addr() string {
	return ""
}

func (ses *emptySession) Send(msg interface{}) error {
	return fmt.Errorf("send message %+v to empty session", msg)
}

func (ses *emptySession) SendStream(stream []byte) error {
	return fmt.Errorf("send stream %s to empty session", string(stream))
}

func (ses *emptySession) Close() {

}

func (ses *emptySession) Set(key interface{}, value interface{}) {
	ses.values.Store(key, value)
}

func (ses *emptySession) Get(key interface{}) (value interface{}, ok bool) {
	return ses.values.Load(key)
}

var globalEmptySession = &emptySession{}

func DefaultEmptySession() Session {
	return globalEmptySession
}
