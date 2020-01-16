package message

import (
	"reflect"
)

type ID uint32

type Meta interface {
	ID() ID
	Name() string
	Codec() Codec
	Encode(msg interface{}) (raw []byte, err error)
	Decode(raw []byte) (msg interface{}, err error)
}

type meta struct {
	id    ID
	name  string
	typo  reflect.Type
	codec Codec
}

func (m *meta) ID() ID {
	return m.id
}

func (m *meta) Name() string {
	return m.name
}

func (m *meta) Codec() Codec {
	return m.codec
}

func (m *meta) Encode(msg interface{}) ([]byte, error) {
	return m.codec.Encode(msg)
}

func (m *meta) Decode(raw []byte) (interface{}, error) {
	msg := reflect.New(m.typo).Interface()
	err := m.codec.Decode(raw, msg)
	return msg, err
}
