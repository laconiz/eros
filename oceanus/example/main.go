package main

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/oceanus/process"
	"github.com/laconiz/eros/utils/command"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	addr := command.ParseStringArg("addr", "")

	proc, err := process.New(addr, encoder.NewNameMaker().New())
	if err != nil {
		panic(err)
	}

	proc.Run()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-exit

	proc.Stop()

	os.Exit(0)
}

var log = logisor.Field(logis.Module, "main")
