package http

import (
	"github.com/gin-gonic/gin"
)

type AcceptorConfig struct {
	Name   string
	Addr   string
	Node   *Node
	Engine *gin.Engine
}

func (c *AcceptorConfig) Load() {

	if c.Name == "" {
		c.Name = "acceptor"
	}
}
