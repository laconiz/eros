package encoder

import (
	"github.com/laconiz/eros/network/message"
)

type Maker interface {
	New() message.Encoder
}
