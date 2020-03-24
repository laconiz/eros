package httpis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/context"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/utils/ioc"
	"net/http"
	"reflect"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

type Node struct {
	Path    string       // 路径
	Method  string       // 方法
	Handler interface{}  // 接口
	Log     bool         // 日志
	typo    reflect.Type // 消息类型
}

// ---------------------------------------------------------------------------------------------------------------------

func NewInvoker(logger logis.Logger) *Invoker {
	return &Invoker{squirt: ioc.New(), logger: logger}
}

// ---------------------------------------------------------------------------------------------------------------------

type Invoker struct {
	squirt *ioc.Squirt  // 依赖注入器
	logger logis.Logger // 日志接口
}

// ---------------------------------------------------------------------------------------------------------------------

func (invoker *Invoker) Params(params ...interface{}) *Invoker {
	invoker.squirt.Params(params...)
	return invoker
}

func (invoker *Invoker) Creators(creators ...interface{}) *Invoker {
	invoker.squirt.Creators(creators...)
	return invoker
}

// ---------------------------------------------------------------------------------------------------------------------

func (invoker *Invoker) Register(router gin.IRouter, nodes []*Node) error {

	for _, node := range nodes {
		if err := invoker.RegisterEx(router, node); err != nil {
			return err
		}
	}

	return nil
}

func (invoker *Invoker) RegisterEx(router gin.IRouter, node *Node) error {

	squirt := invoker.squirt.Copy()

	args, err := squirt.UnknownArgs(node.Handler, invoker.dynamicParams()...)
	if err != nil {
		return err
	}

	if len(args) > 1 {
		return fmt.Errorf("too many args in handler: %v", args)
	}

	if len(args) == 1 {

		typo := args[0]

		if typo.Kind() != reflect.Ptr {
			return fmt.Errorf("message type must be pointer: %v", typo)
		}

		squirt.Creator(typo, func(ctx *gin.Context) (interface{}, error) {
			msg := reflect.New(typo.Elem()).Interface()
			err := ctx.Bind(msg)
			return msg, err
		})
	}

	handler, err := squirt.Handle(node.Handler, invoker.dynamicParams()...)
	if err != nil {
		return err
	}

	invoker.logger.Infof("%s: %s", node.Method, node.Path)
	router.Handle(node.Method, node.Path, invoker.Handle(node, handler))
	return nil
}

func (invoker *Invoker) Handle(node *Node, handler ioc.Invoker) gin.HandlerFunc {

	log := invoker.logger.Fields(context.Fields{fieldPath: node.Path, fieldMethod: node.Method})

	return func(ctx *gin.Context) {

		log = log.Field(network.FieldSession, session.Increment())

		now := time.Now()
		args, values, err := handler(ctx)
		log = log.Field(fieldDuration, time.Since(now).Milliseconds())

		if len(args) > 0 {
			log = log.Field(fieldRequests, args)
		}

		if err != nil {
			log.Warnf("invoke error: %v", err)
			return
		}

		var responses []interface{}
		for _, value := range values {
			if !value.CanInterface() {
				continue
			}
			responses = append(responses, value.Interface())
		}

		for _, response := range responses {
			ctx.JSON(http.StatusOK, response)
		}
		log = log.Field(fieldResponses, responses)

		log.Info("execute success")
	}
}

func (invoker *Invoker) dynamicParams() []interface{} {
	return []interface{}{&gin.Context{}}
}

// ---------------------------------------------------------------------------------------------------------------------

const (
	fieldPath      = "path"
	fieldMethod    = "method"
	fieldDuration  = "milliseconds"
	fieldRequests  = "requests"
	fieldResponses = "responses"
)
