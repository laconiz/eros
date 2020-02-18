package main

import (
	"github.com/gin-gonic/gin"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network/invoker"
	"net/http"
	"os"
	"time"

	"github.com/laconiz/eros/network/httpis"
	"github.com/laconiz/eros/network/httpis/example"
)

func main() {

	var state string

	log := logis.NewHook(logis.NewTextFormatter()).AddWriter(logis.DEBUG, os.Stdout).Entry()

	opt := httpis.AcceptorOption{
		Addr: example.Addr,
		Nodes: []*invoker.Node{
			{
				Path:   example.Path,
				Method: http.MethodGet,
				Handler: func() *example.StateACK {
					return &example.StateACK{State: state, Time: time.Now()}
				},
			},
			{
				Path:   example.Path,
				Method: http.MethodPut,
				Handler: func(req *example.ReportREQ, ctx *gin.Context) string {
					state = req.State
					return state
				},
			},
			{
				Path:   example.Path,
				Method: http.MethodPost,
				Handler: func(req *example.ReportREQ) bool {
					state = req.State
					return true
				},
			},
		},
	}

	acceptor, err := httpis.NewAcceptor(opt, log)
	if err != nil {
		panic(err)
	}
	acceptor.Run()

	ch := make(chan bool)
	<-ch
}
