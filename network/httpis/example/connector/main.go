package main

import (
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network/httpis"
	"github.com/laconiz/eros/network/httpis/example"
)

func main() {

	connector := httpis.URL(example.Addr + example.Path)

	var state string
	if err := connector.Put(&example.ReportREQ{State: "PUT"}, &state); err != nil {
		panic(err)
	}
	logger.Info(state)

	ack := example.StateACK{}
	if err := connector.Get(nil, &ack); err != nil {
		panic(err)
	}
	logger.Info(ack)

	var success bool
	if err := connector.Post(&example.ReportREQ{State: "normal"}, &success); err != nil {
		panic(err)
	}
	logger.Info(success)

	if err := connector.Get(nil, &ack); err != nil {
		panic(err)
	}
	logger.Info(ack)
}

var logger = logisor.Module("main")
