package main

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network/steropes"
	"github.com/laconiz/eros/network/steropes/example"
)

func main() {

	connector := steropes.URL(example.Addr + example.Path)

	if err := connector.Post(&example.ReportREQ{State: "normal"}, nil); err != nil {
		panic(err)
	}

	ack := &example.StateACK{}
	if err := connector.Get(nil, ack); err != nil {
		panic(err)
	}

	logger.Info(*ack)
}

var logger = logis.NewEntry("main")
