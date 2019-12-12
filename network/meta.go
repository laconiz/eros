package network

import (
	"reflect"
)

type MessageID uint32

type Meta struct {
	id    MessageID
	typo  reflect.Type
	codec Codec
}

func (m *Meta) ID() MessageID {
	return m.id
}

func (m *Meta) Name() string {
	return m.typo.String()
}

func (m *Meta) Type() reflect.Type {
	return m.typo
}

func (m *Meta) Codec() Codec {
	return m.codec
}

func (m *Meta) String() string {
	return m.typo.String()
}
