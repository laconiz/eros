package network

import (
	"github.com/laconiz/eros/json"
)

type Codec interface {
	Encode(msg interface{}) (raw []byte, err error)
	Decode(raw []byte, msg interface{}) (err error)
}

type jsonCodec struct {
}

func (c *jsonCodec) Encode(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

func (c *jsonCodec) Decode(raw []byte, msg interface{}) error {
	return json.Unmarshal(raw, msg)
}

var JsonCodec Codec = &jsonCodec{}
