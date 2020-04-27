package reader

import "net"

// ---------------------------------------------------------------------------------------------------------------------

type Reader interface {
	Write(net.Conn, []byte) error
	Read(net.Conn) ([]byte, error)
}

// ---------------------------------------------------------------------------------------------------------------------

type Maker interface {
	New() Reader
}