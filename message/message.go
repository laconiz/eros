package message

import (
	"fmt"
	"github.com/laconiz/eros/json"
)

type Message interface {
	Meta() Meta
	Message() interface{}
	Raw() []byte
	Stream() []byte
	String() string
}

type message struct {
	meta    Meta
	msg     interface{}
	raw     []byte
	stream  []byte
	encoder Encoder
}

func (m *message) Meta() Meta {
	return m.meta
}

func (m *message) Message() interface{} {
	return m.msg
}

func (m *message) Raw() []byte {
	return m.raw
}

func (m *message) Stream() []byte {
	return m.stream
}

func (m *message) String() string {

	if m.meta.Codec() == globalJsonCodec {

		if m.encoder == globalNameEncoder {
			return string(m.stream)
		}

		return fmt.Sprintf("%s-%s", m.meta.Name(), string(m.raw))
	}

	raw, err := json.Marshal(m.msg)
	if err != nil {
		return fmt.Sprintf("%s-%v", m.meta.Name(), err)
	}
	return fmt.Sprintf("%s-%s", m.meta.Name(), string(raw))
}
