package main

import (
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/example"
	"math/rand"
	"time"
)

var times time.Duration
var duration time.Duration

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewConnector() *socket.Connector {

	invoker := invoker.NewSocketInvoker()
	invoker.Register(example.ACK{}, func(event *network.Event) {
		ack := event.Msg.(*example.ACK)
		times++
		duration += time.Since(ack.Time)
	})

	option := &socket.ConnectorOption{
		Addr:      example.Addr,
		Reconnect: true,
		Session: socket.SessionOption{
			Invoker: invoker,
		},
	}

	return socket.NewConnector(option)
}

func main() {

	connector := NewConnector()
	connector.Run()

	go func() {

		var bytes []byte
		for i := 0; i < 4096; i++ {
			bytes = append(bytes, byte(i))
		}

		for {

			connector.Send(&example.REQ{
				Time:  time.Now(),
				Bytes: bytes[:random.Intn(4096)+1],
			})
		}
	}()

	for {
		<-time.After(time.Second)
		if connector.Connected() {
			logger.Infof("%v / %d = %v", duration, times, duration/times)
			times = 0
			duration = 0
		}
	}
}

var logger = logisor.Module("socket.example.connector")
