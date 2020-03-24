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
	Load(key interface{}) (value interface{}, ok bool)
	Store(key interface{}, value interface{})
}

// ---------------------------------------------------------------------------------------------------------------------

type EmptySession struct {
	sync.Map
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
