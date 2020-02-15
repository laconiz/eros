package steropes

import (
	"fmt"
	"github.com/laconiz/eros/network/epimetheus"
	"github.com/laconiz/eros/network/incremental"
	"github.com/laconiz/eros/utils/ioc"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/laconiz/eros/logis"
)

// 节点信息
type Node struct {
	// 路径
	Path string
	// 接口
	Handlers map[string]interface{}
	// 子节点
	Children []*Node
	// 是否记录日志
	NoLog bool
}

func handleNode(router gin.IRouter, node *Node, base *ioc.Squirt, logger *logis.Entry) error {

	for method, handler := range node.Handlers {

		// 获取未知参数
		args, err := base.UnknownArgs(handler)
		if err != nil {
			return err
		}
		// 消息参数数量过多
		if len(args) > 1 {
			return fmt.Errorf("too much messages: %v", args)
		}

		squirt := base
		// 更新依赖注入器
		if len(args) == 1 {

			// 检查消息类型
			typo := args[0]
			if typo.Kind() != reflect.Ptr || typo.Elem().Kind() != reflect.Struct {
				return fmt.Errorf("message type must be pointer of struct, got %v", typo)
			}

			// 注入消息生成器
			squirt = base.Copy().Function(typo, func(ctx *gin.Context) (interface{}, error) {
				message := reflect.New(typo.Elem()).Interface()
				err := ctx.Bind(message)
				return message, err
			})
		}

		// 生成调用接口
		invoker, err := squirt.Handle(handler, &gin.Context{})
		if err != nil {
			return err
		}

		path := router.(*gin.RouterGroup).BasePath()
		path = strings.Replace(path+node.Path, "//", "/", -1)
		logger = logger.WithField(fieldMethod, method).
			WithField(fieldPath, path+node.Path)

		logger.Info("registered")
		router.Handle(method, node.Path, handleRouter(invoker, handleResponse(handler), logger, !node.NoLog))
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

func handleRouter(invoker ioc.Invoker, handler ResponseHandler, logger *logis.Entry, log bool) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		now := time.Now()
		logger = logger.WithField(epimetheus.FieldSession, incremental.Get())

		// 调用接口
		arguments, values, err := invoker(ctx)
		if err != nil {
			logger.WithField(fieldDuration, time.Since(now)/time.Nanosecond).
				WithError(err).Error("invoke error")
			logger.Info(err.Error())
			return
		}
		if arguments != nil {
			logger = logger.WithField(fieldArgument, arguments)
		}

		// 调用返回值处理接口
		responses, err := handler(ctx, values)
		if err != nil {
			logger.WithField(fieldDuration, time.Since(now)/time.Nanosecond).
				WithError(err).Error("handle result error")
			return
		}
		if responses != nil {
			logger = logger.WithField(fieldResponse, responses)
		}

		if log {
			logger.WithField(fieldDuration, time.Since(now)/time.Nanosecond).
				Info("execute success")
		}
	}
}

type ResponseHandler func(*gin.Context, []reflect.Value) ([]interface{}, error)

// 构建返回值响应函数
func handleResponse(handler interface{}) ResponseHandler {

	typo := reflect.TypeOf(handler)

	if typo.NumOut() == 0 {
		return func(_ *gin.Context, _ []reflect.Value) ([]interface{}, error) {
			return nil, nil
		}
	}

	lastOut := typo.Out(typo.NumOut() - 1)
	if lastOut.Implements(reflect.TypeOf((*error)(nil)).Elem()) {

		return func(ctx *gin.Context, values []reflect.Value) ([]interface{}, error) {
			err := values[len(values)-1]
			if !err.IsNil() {
				return nil, err.Interface().(error)
			}
			var responses []interface{}
			for i := 0; i < len(values)-1; i++ {
				if values[i].CanInterface() {
					response := values[i].Interface()
					ctx.JSON(http.StatusOK, response)
					responses = append(responses, response)
				}
			}
			return responses, nil
		}

	} else {

		return func(ctx *gin.Context, values []reflect.Value) ([]interface{}, error) {
			var responses []interface{}
			for _, value := range values {
				if value.CanInterface() {
					response := value.Interface()
					ctx.JSON(http.StatusOK, response)
					responses = append(responses, response)
				}
			}
			return responses, nil
		}
	}
}

const (
	fieldPath     = "path"
	fieldMethod   = "method"
	fieldArgument = "argument"
	fieldResponse = "response"
	fieldDuration = "duration"
)
