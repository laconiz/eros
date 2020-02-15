package session

import (
	"fmt"
	"math"
	"sync"
)

type ID uint64

type Session interface {
	ID() ID
	Addr() string
	Send(msg interface{}) error
	SendRaw(raw []byte) error
	Close()
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
}

// ---------------------------------------------------------------------------------------------------------------------

type EmptySession struct {
	values sync.Map
}

func (ses *EmptySession) ID() ID {
	return math.MaxUint64
}

func (ses *EmptySession) Addr() string {
	return ""
}

func (ses *EmptySession) Send(msg interface{}) error {
	return fmt.Errorf("send message %+v to empty session", msg)
}

func (ses *EmptySession) SendRaw(raw []byte) error {
	return fmt.Errorf("send stream %s to empty session", string(raw))
}

func (ses *EmptySession) Close() {

}

func (ses *EmptySession) Set(key interface{}, value interface{}) {
	ses.values.Store(key, value)
}

func (ses *EmptySession) Get(key interface{}) (interface{}, bool) {
	return ses.values.Load(key)
}
