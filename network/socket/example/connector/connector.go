package main

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/formatter"
	"github.com/laconiz/eros/logis/hook"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/example"
	"os"
)

func main() {

	opt := socket.ConnectorOption{
		Addr:      example.Addr,
		Reconnect: true,
	}
	conn := socket.NewConnector(opt)
	conn.Run()

	for {
		conn.Send(&example.REQ{Int: 1})
	}
}

var log = hook.NewHook(formatter.Text()).Add(logis.DEBUG, os.Stdout).Entry()
