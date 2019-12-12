package http

import (
	"github.com/gin-gonic/gin"
	"github.com/laconiz/eros/network/invoker"
	"strings"
)

type Node struct {
	Path     string
	Handlers map[string]interface{}
	Children []*Node
	NoLog    bool
}

func HandleNode(router gin.IRouter, node *Node, name string, params ...interface{}) error {

	for method, handler := range node.Handlers {

		path := router.(*gin.RouterGroup).BasePath()
		path = strings.Replace(path+node.Path, "//", "/", -1)

		invoker, err := invoker.NewHttpInvoker(name, method, path, handler, params...)
		if err != nil {
			return err
		}

		router.Handle(method, path, invoker)
	}

	group := router.Group(node.Path)
	for _, child := range node.Children {
		if err := HandleNode(group, child, name, params...); err != nil {
			return err
		}
	}

	return nil
}
