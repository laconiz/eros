package main

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/network/websocket"
	"github.com/laconiz/eros/network/websocket/examples"
)

func onREQ(req *examples.REQ) *examples.ACK {
	// time.Sleep(time.Second)
	return &examples.ACK{ID: req.ID + 1}
}

func main() {

	logger := log.Std("main")

	handlers := []interface{}{onREQ}

	inv, err := invoker.NewNetworkInvoker(logger, handlers)
	if err != nil {
		panic(err)
	}

	conf := websocket.AcceptorConfig{
		Addr: ":12314",
	}
	conf.Session.Invoker = inv

	acc := websocket.NewAcceptor(conf)
	acc.Start()

	c := make(chan bool)
	<-c
}
