package main

import (
	"github.com/laconiz/eros/oceanus"
	"os"
)

func main() {

	process := oceanus.NewProcess(os.Args[1])
	process.Run()

}
