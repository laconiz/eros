package main

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/examples"
)

func main() {

	opt := socket.ConnOption{
		Addr:      "127.0.0.1:1024",
		Reconnect: true,
	}
	conn := socket.NewConnector(opt)
	conn.Run()

	for {
		conn.Send(&examples.REQ{Int: 1})
	}
}

var log = logisor.Field(logis.Module, "main")
