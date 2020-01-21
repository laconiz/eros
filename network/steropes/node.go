package steropes

import (
	"net/http"
	"reflect"

	"github.com/codegangsta/inject"
	"github.com/gin-gonic/gin"

	"github.com/laconiz/eros/hyperion"
)

type Node struct {
	Path     string
	Handlers map[string]interface{}
	Children []*Node
	NoLog    bool
}

func newHandler(node *Node) (http.Handler, error) {

	engine := gin.New()
	engine.Use(gin.Recovery())

	return engine, nil
}

type getter struct {
	typo    reflect.Type
	handler interface{}
}

type invoker struct {
	logger   *hyperion.Entry
	method   string
	path     string
	injector inject.Injector
	getters  []*getter
}
