package main

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/network/websocket"
	"github.com/laconiz/eros/network/websocket/examples"
)

func onACK(ack *examples.ACK, flag *uint) *examples.REQ {
	// time.Sleep(time.Second)
	*flag = ack.ID + 1
	return &examples.REQ{ID: *flag}
}

func onConnected(_ *network.Connected, flag *uint) *examples.REQ {
	return &examples.REQ{ID: *flag}
}

func onConnectedFailed(_ *network.ConnectFailed) {
	logger.Warn("connect failed")
}

func main() {

	flag := uint(0)

	handlers := []interface{}{onConnected, onConnectedFailed, onACK}

	inv, err := invoker.NewNetworkInvoker(logger, handlers, &flag)
	if err != nil {
		panic(err)
	}

	conf := websocket.ConnectorConfig{
		Addr:      "ws://127.0.0.1:12314/ws",
		Reconnect: true,
		Session: websocket.SessionConfig{
			Invoker: inv,
		},
	}

	cnt := websocket.NewConnector(conf)
	cnt.Start()

	c := make(chan bool)
	<-c
}

var logger = log.Std("main")
