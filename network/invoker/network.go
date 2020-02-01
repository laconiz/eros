package invoker

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/codegangsta/inject"

	"github.com/laconiz/eros/holder/message"
	"github.com/laconiz/eros/network"
)

type SocketInvoker struct {
	log      Logger                               // 日志接口
	handlers map[message.ID][]network.HandlerFunc // 消息接口
}

// 调用接口
func (inv *SocketInvoker) Invoke(event *network.Event) {
	for _, handlers := range inv.handlers[event.ID] {
		func() {
			defer func() {
				if err := recover(); err != nil {
					inv.log.Errorf("capture panic: %v", err)
				}
			}()
			handlers(event)
		}()
	}
}

// 构造消息回调器
// logger	日志接口
// handlers	消息接口
// params	基础参数和参数生成器
func NewNetworkInvoker(log Logger, handlers []interface{}, params ...interface{}) (network.Invoker, error) {

	inv := &NetworkInvoker{
		log:      log,
		handlers: map[network.MessageID][]network.HandlerFunc{},
	}

	for _, handler := range handlers {

		meta, invoker, err := newNetworkHandler(log, handler, params...)
		if err != nil {
			return nil, err
		}

		inv.handlers[meta.ID()] = append(inv.handlers[meta.ID()], func(event *network.Event) {

			// 防止上层逻辑导致的崩溃
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("invoke message[%+v] error: %v", event.Msg, err)
				}
			}()

			invoker(event)
		})
	}

	return inv, nil
}

// 构造消息回调
// logger	日志接口
// handler	消息接口
// params	基础参数和参数生成器
func newNetworkHandler(
	log Logger,
	handler interface{},
	params ...interface{},
) (*network.Meta, network.HandlerFunc, error) {

	// 检查接口类型
	typo := reflect.TypeOf(handler)
	if typo == nil || typo.Kind() != reflect.Func {
		return nil, nil, fmt.Errorf("need func, got %v", typo)
	}

	// 生成基础调用信息
	invoker, err := newInvoker(params...)
	if err != nil {
		return nil, nil, fmt.Errorf("parse param error: %w", err)
	}

	// 生成基础调用链
	if err := invoker.makeChains(&network.Event{}, network.DefaultSession); err != nil {
		return nil, nil, fmt.Errorf("make chain error: %w", err)
	}

	// 接口返回值检查
	for i := 0; i < typo.NumOut(); i++ {
		if network.MetaByType(typo.Out(i)) == nil {
			return nil, nil, fmt.Errorf("invalid out: %v", typo.Out(i))
		}
	}

	// 消息接口调用信息
	inv := &networkHandler{
		Injector: invoker.Injector,
		Handler:  handler,
		log:      log,
	}

	// 可提供给参数生成器调用的参数
	injector := inject.New()
	injector.SetParent(inv.Injector)
	injector.Map(&network.Event{}).Map(network.Session(nil))

	// 消息元数据
	var meta *network.Meta

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

		// 多个消息参数
		if meta != nil {
			return nil, nil, fmt.Errorf("invalid argument: %v", typo.In(i))
		}

		// 未注册的消息参数
		if meta = network.MetaByType(typo.In(i)); meta == nil {
			return nil, nil, fmt.Errorf("invalid argument: %v", typo.In(i))
		}
	}

	if meta == nil {
		return nil, nil, errors.New("message not found")
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

	return meta, inv.Invoke, nil
}

// 消息接口调用信息
type networkHandler struct {
	log      Logger          // 日志接口
	Injector inject.Injector // 基础参数
	Getters  []*paramGetter  // 参数生成链
	Handler  interface{}     // 调用接口
}

// 调用接口
func (inv *networkHandler) Invoke(event *network.Event) {

	// 基础参数列表
	inj := inject.New()
	inj.SetParent(inv.Injector)
	inj.Map(event).Map(event.Session).Map(event.Msg)

	// 调用参数生成器
	for _, getter := range inv.Getters {

		// 调用失败
		values, err := inj.Invoke(getter.Handler)
		if err != nil {
			inv.log.Errorf("invoke getter %v error: %v", getter.Type, err)
			return
		}

		// 获取参数失败
		if !values[1].IsNil() {
			inv.log.Warnf("get param %v error: %v", getter.Type, err)
			return
		}

		// 写入参数
		inj.Set(getter.Type, values[0])
	}

	// 调用接口
	values, err := inj.Invoke(inv.Handler)
	if err != nil {
		inv.log.Errorf("invoke handler error: %v", err)
		return
	}

	// 发送返回值
	if event.Session != nil {
		for _, value := range values {
			event.Session.Send(value.Interface())
		}
	}
}
