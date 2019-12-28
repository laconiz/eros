package codec

import "github.com/laconiz/eros/json"

type jsonCodec struct {
}

func (c *jsonCodec) Encode(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

func (c *jsonCodec) Decode(raw []byte, msg interface{}) error {
	return json.Unmarshal(raw, msg)
}

var globalJsonCodec = &jsonCodec{}

func Json() Codec {
	return globalJsonCodec
}
