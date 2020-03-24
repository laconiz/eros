package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/laconiz/eros/network/httpis"
	"github.com/laconiz/eros/network/httpis/example"
)

func main() {

	var state string

	option := &httpis.AcceptorOption{
		Addr: example.Addr,
		Nodes: []*httpis.Node{
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

	acceptor, err := httpis.NewAcceptor(option)
	if err != nil {
		panic(err)
	}
	acceptor.Run()

	ch := make(chan bool)
	<-ch
}
