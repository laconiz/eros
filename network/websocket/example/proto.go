package example

import (
	"github.com/laconiz/eros/network/message"
	"time"
)

const Addr = "192.168.1.2:8001"

type REQ struct {
	Time  time.Time
	Bytes []byte
}

type ACK struct {
	Time time.Time
}

func init() {
	message.Json(REQ{})
	message.Json(ACK{})
}
