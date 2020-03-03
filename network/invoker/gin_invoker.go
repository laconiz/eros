package invoker

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

type Node struct {
	Path    string       // 路径
	Method  string       // 方法
	Handler interface{}  // 接口
	Log     bool         // 日志
	typo    reflect.Type // 消息类型
}

// ---------------------------------------------------------------------------------------------------------------------

func NewGinInvoker(log logis.Logger) *GinInvoker {
	return &GinInvoker{squirt: ioc.New(), log: log}
}

type GinInvoker struct {
	squirt *ioc.Squirt
	log    logis.Logger
}

func (i *GinInvoker) Params(params ...interface{}) *GinInvoker {
	i.squirt.Params(params...)
	return i
}

func (i *GinInvoker) Creators(creators ...interface{}) *GinInvoker {
	i.squirt.Creators(creators...)
	return i
}

func (i *GinInvoker) Register(router gin.IRouter, nodes []*Node) error {

	for _, node := range nodes {
		if err := i.RegisterEx(router, node); err != nil {
			return err
		}
	}

	return nil
}

func (i *GinInvoker) RegisterEx(router gin.IRouter, node *Node) error {

	squirt := i.squirt.Copy()

	args, err := squirt.UnknownArgs(node.Handler, i.dynamicParams()...)
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

	invoker, err := squirt.Handle(node.Handler, i.dynamicParams()...)
	if err != nil {
		return err
	}

	i.log.Infof("%s: %s", node.Method, node.Path)
	router.Handle(node.Method, node.Path, i.Handle(node, invoker))
	return nil
}

// 将依赖注入接口转换为gin接口
func (i *GinInvoker) Handle(node *Node, invoker ioc.Invoker) gin.HandlerFunc {

	log := i.log.Fields(context.Fields{fieldPath: node.Path, fieldMethod: node.Method})

	return func(ctx *gin.Context) {

		log = log.Field(network.FieldSession, session.Increment())

		now := time.Now()
		args, values, err := invoker(ctx)
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

func (i *GinInvoker) dynamicParams() []interface{} {
	return []interface{}{&gin.Context{}}
}

const (
	fieldPath      = "path"
	fieldMethod    = "method"
	fieldDuration  = "duration"
	fieldRequests  = "requests"
	fieldResponses = "responses"
)
