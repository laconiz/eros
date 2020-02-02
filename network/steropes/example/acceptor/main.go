package main

import (
	"github.com/laconiz/eros/network/steropes"
	"github.com/laconiz/eros/network/steropes/example"
	"net/http"
	"time"
)

func main() {

	var state string

	acceptor, err := steropes.NewAcceptor(steropes.AcceptorOption{
		Addr: example.Addr,
		Node: &steropes.Node{
			Path:     "/",
			Handlers: nil,
			Children: []*steropes.Node{
				{
					Path: example.Path,
					Handlers: map[string]interface{}{
						http.MethodPost: func(req *example.ReportREQ) bool {
							state = req.State
							return true
						},
						http.MethodGet: func() *example.StateACK {
							return &example.StateACK{State: state, Time: time.Now()}
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	acceptor.Run()
}
