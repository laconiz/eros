package encoder

import (
	"github.com/laconiz/eros/holder/message"
)

type Maker interface {
	New() message.Encoder
}
