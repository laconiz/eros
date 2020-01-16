// 消息序列化规则

package message

import "github.com/laconiz/eros/json"

type Codec interface {
	Encode(msg interface{}) (raw []byte, err error)
	Decode(raw []byte, msg interface{}) (err error)
}

// ---------------------------------------------------------------------------------------------------------------------
// JSON序列化规则

type jsonCodec struct {
}

func (c *jsonCodec) Encode(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

func (c *jsonCodec) Decode(raw []byte, msg interface{}) error {
	return json.Unmarshal(raw, msg)
}

var globalJsonCodec = &jsonCodec{}

func JsonCodec() Codec {
	return globalJsonCodec
}
