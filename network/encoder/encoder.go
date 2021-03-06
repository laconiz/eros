package encoder

import "github.com/laconiz/eros/network/message"

// ---------------------------------------------------------------------------------------------------------------------

type Encoder interface {
	Marshal(msg interface{}) (*message.Message, error)
	Unmarshal(stream []byte) (*message.Message, error)
}

// ---------------------------------------------------------------------------------------------------------------------

type Maker interface {
	New() Encoder
}
