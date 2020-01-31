package steropes

import (
	"fmt"
	"github.com/laconiz/eros/utils/ioc"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/laconiz/eros/hyperion"
)

type Node struct {
	Path     string
	Handlers map[string]interface{}
	Children []*Node
	NoLog    bool
}

func newHandler(node *Node, logger *hyperion.Entry) (http.Handler, error) {
	// gin
	engine := gin.New()
	engine.Use(gin.Recovery())
	// 注入器
	squirt := ioc.New().Params().Functions()
	if err := squirt.Error(); err != nil {
		return nil, err
	}
	// 构造接口
	if err := handleNode(engine, node, squirt, logger); err != nil {
		return nil, err
	}
	return engine, nil
}

func handleNode(router gin.IRouter, node *Node, base *ioc.Squirt, logger *hyperion.Entry) error {

	for method, handler := range node.Handlers {

		squirt, err := addMessageHandler(base, handler)
		if err != nil {
			return err
		}

		invoker, err := squirt.Handle(handler, &gin.Context{})
		if err != nil {
			return err
		}

		path := router.(*gin.RouterGroup).BasePath()
		logger = logger.WithField(fieldPath, path+node.Path)

		router.Handle(method, node.Path, func(ctx *gin.Context) {
			invoker(ctx)
		})
	}
	// 响应子路径
	router = router.Group(node.Path)
	for _, child := range node.Children {
		if err := handleNode(router, child, base, logger); err != nil {
			return err
		}
	}
	return nil
}

// 构造消息生成器
func addMessageHandler(squirt *ioc.Squirt, handler interface{}) (*ioc.Squirt, error) {
	// 获取未知参数列表
	args, err := squirt.UnknownArgs(handler)
	if err != nil {
		return nil, err
	}
	// 未知参数过多
	if len(args) > 1 {
		return nil, fmt.Errorf("too many unknown args %v", args)
	}
	if len(args) == 1 {
		squirt = squirt.Copy()
		// 检查消息类型
		typo := args[0]
		if typo.Kind() != reflect.Ptr {
			return nil, invalidMessageTypeError(typo)
		}
		typo = typo.Elem()
		if typo.Kind() != reflect.Struct {
			return nil, invalidMessageTypeError(typo)
		}
		// 插入消息生成器
		squirt.Function(typo, func(ctx *gin.Context) (interface{}, error) {
			message := reflect.New(typo).Interface()
			err := ctx.Bind(message)
			return message, err
		})
	}
	return squirt, nil
}

func invalidMessageTypeError(typo reflect.Type) error {
	return fmt.Errorf("message type must be pointer of struct, got %v", typo)
}

const (
	fieldPath = "path"
)
