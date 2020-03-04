package main

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/oceanus/process"
	"github.com/laconiz/eros/utils/command"
)

func main() {

	addr := command.ParseStringArg("addr", "")
	log.Info(addr)

	proc, err := process.New(addr)
	if err != nil {
		panic(err)
	}

	proc.Run()
}

var log = logisor.Field(logis.Module, "main")
