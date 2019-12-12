package invoker

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"sync/atomic"
)

// 构造消息回调
// logger	日志接口
// method	方法
// path		路径
// handler	消息接口
// params	基础参数和参数生成器
func NewHttpInvoker(
	logger Logger,
	method string,
	path string,
	handler interface{},
	params ...interface{},
) (gin.HandlerFunc, error) {

	// 检查接口类型
	typo := reflect.TypeOf(handler)
	if typo == nil || typo.Kind() != reflect.Func {
		return nil, fmt.Errorf("need func, got %v", typo)
	}

	// 生成基础调用信息
	invoker, err := newInvoker(params...)
	if err != nil {
		return nil, fmt.Errorf("parse param error: %w", err)
	}

	// 生成基础调用链
	if err := invoker.makeChains(&gin.Context{}); err != nil {
		return nil, fmt.Errorf("make chain error: %w", err)
	}

	// 消息接口调用信息
	inv := &httpInvoker{
		Logger:   logger,
		Method:   method,
		Path:     path,
		Injector: invoker.Injector,
		Handler:  handler,
	}

	// 可提供给参数生成器调用的参数
	injector := inject.New()
	injector.SetParent(inv.Injector)
	injector.Map(&gin.Context{})

	// 消息类型
	var message reflect.Type

	for i := 0; i < typo.NumIn(); i++ {

		// 基础参数
		if injector.Get(typo.In(i)).IsValid() {
			continue
		}

		// 参数生成器
		if getters, ok := invoker.Chains[typo.In(i)]; ok {
			inv.Getters = append(inv.Getters, getters...)
			continue
		}

		// 消息生成器
		if message == nil {
			message = typo.In(i)
			inv.Getters = append(inv.Getters, makeMessageChain(message))
			continue
		}

		// 多个消息参数
		return nil, fmt.Errorf("invalid argument: %v", typo.In(i))
	}

	// 参数生成器去重
	var getters []*paramGetter
	for _, getter := range inv.Getters {
		for i := 0; i < len(getters); i++ {
			if getters[i].Type == getter.Type {
				goto SKIP
			}
		}
		getters = append(getters, getter)
	SKIP:
	}
	inv.Getters = getters

	return inv.Invoke, nil
}

// 构造消息生成器
func makeMessageChain(typo reflect.Type) *paramGetter {

	if typo.Kind() == reflect.Ptr {

		// 获取消息指针
		return &paramGetter{
			Type: typo,
			Handler: func(ctx *gin.Context) (interface{}, error) {
				message := reflect.New(typo.Elem()).Interface()
				err := ctx.Bind(message)
				return message, err
			},
		}
	}

	// 获取消息值
	return &paramGetter{
		Type: typo,
		Handler: func(ctx *gin.Context) (interface{}, error) {
			message := reflect.New(typo)
			err := ctx.Bind(message.Interface())
			return message.Elem().Interface(), err
		},
	}
}

// 消息接口调用信息
type httpInvoker struct {
	Logger   Logger          // 日志接口
	Method   string          // 调用方式
	Path     string          // 调用路径
	Injector inject.Injector // 基础参数
	Getters  []*paramGetter  // 参数生成链
	Handler  interface{}     // 调用接口
}

var requestID uint64 = 0

func (inv *httpInvoker) Invoke(ctx *gin.Context) {

	// 全局唯一标识
	id := atomic.AddUint64(&requestID, 1)
	inv.Logger.Infof("request[%d] %s %s", id, inv.Method, inv.Path)

	// 基础参数列表
	inj := inject.New()
	inj.SetParent(inv.Injector)
	inj.Map(ctx)

	// 调用参数生成器
	for _, getter := range inv.Getters {

		// 调用失败
		values, err := inj.Invoke(getter.Handler)
		if err != nil {
			inv.Logger.Errorf("request[%d] invoke getter %v error: %v", id, getter.Type, err)
			return
		}

		// 获取参数失败
		if !values[1].IsNil() {
			inv.Logger.Warnf("request[%d] get param %v error: %v", id, getter.Type, err)
			return
		}

		// 记录日志
		raw, err := json.Marshal(values[0].Interface())
		inv.Logger.Infof("request[%d] param: %v<%s> error: %v", id, values[0].Type(), string(raw), err)

		// 写入参数
		inj.Set(getter.Type, values[0])
	}

	// 调用接口
	values, err := inj.Invoke(inv.Handler)
	if err != nil {
		inv.Logger.Errorf("request[%d] invoke handler error: %v", id, err)
		return
	}

	for _, value := range values {

		// 记录日志
		raw, err := json.Marshal(value.Interface())
		inv.Logger.Infof("request[%d] response: %v<%s> error: %v", id, value.Type(), string(raw), err)

		// 写入返回值
		ctx.JSON(http.StatusOK, value.Interface())
	}

	inv.Logger.Infof("request[%d] complete", id)
}
