package encoder

import "github.com/laconiz/eros/network/message"

// ---------------------------------------------------------------------------------------------------------------------

type Encoder interface {
	Encode(msg interface{}) (*message.Message, error)
	Decode(stream []byte) (*message.Message, error)
}

// ---------------------------------------------------------------------------------------------------------------------

type Maker interface {
	New() Encoder
}
