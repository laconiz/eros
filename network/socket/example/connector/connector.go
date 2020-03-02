package main

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/example"
)

func main() {

	opt := socket.ConnOption{
		Addr:      example.Addr,
		Reconnect: true,
	}
	conn := socket.NewConnector(opt)
	conn.Run()

	for {
		conn.Send(&example.REQ{Int: 1})
	}
}

var log = logis.NewHook(logis.NewTextFormatter()).Entry()
