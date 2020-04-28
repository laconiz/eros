package main

import (
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/oceanus/process"
	"github.com/laconiz/eros/utils/command"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	const arg = "addr"
	const value = ""

	// 获取参数
	addr, err := command.ParseAddress(arg, value)
	if err != nil {
		panic(err)
	}

	// 消息序列器
	encoder := encoder.NewNameMaker().New()

	// 创建进程
	proc, err := process.New(addr, encoder)
	if err != nil {
		panic(err)
	}

	// 运行进程
	proc.Run()

	// 响应退出信号
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-exit

	// 退出进程
	proc.Stop()

	os.Exit(0)
}
