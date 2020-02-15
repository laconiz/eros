package socket

import (
	"github.com/laconiz/eros/network/socket/packer"
	"net"
)

func newConn(conn net.Conn, packer packer.Packer) *connection {
	return &connection{Conn: conn, packer: packer}
}

type connection struct {
	net.Conn
	packer packer.Packer
}

func (c *connection) Addr() string {
	return c.Conn.RemoteAddr().String()
}

func (c *connection) Read() ([]byte, error) {
	return c.packer.Decode(c.Conn)
}

func (c *connection) Write(stream []byte) error {
	return c.packer.Encode(c.Conn, stream)
}
