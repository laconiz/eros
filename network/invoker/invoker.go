package invoker

import (
	"errors"
	"fmt"
	"github.com/codegangsta/inject"
	"reflect"
)

func newInvoker(params ...interface{}) (*invoker, error) {

	p := &invoker{
		Injector: inject.New(),
		Handlers: map[reflect.Type]interface{}{},
		Chains:   map[reflect.Type][]*paramGetter{},
	}

	for _, param := range params {

		// 检查参数类型
		typo := reflect.TypeOf(param)
		if typo == nil {
			return nil, errors.New("nil param")
		}

		// 写入基础参数
		if typo.Kind() != reflect.Func {
			// 重复参数
			if p.Injector.Get(typo).IsValid() {
				return nil, fmt.Errorf("conflict type: %v", typo)
			}
			p.Injector.Map(param)
			continue
		}

		// 检查参数生成器
		if typo.NumOut() != 2 || typo.Out(1).Name() != "error" {
			return nil, fmt.Errorf("invalid getter: %v", param)
		}
		// 重复参数
		out := typo.Out(0)
		if p.Injector.Get(out).IsValid() || p.Handlers[out] != nil {
			return nil, fmt.Errorf("conflict type: %v", out)
		}
		// 写入参数生成器
		p.Handlers[out] = param
	}

	return p, nil
}

type paramGetter struct {
	Type    reflect.Type
	Handler interface{}
}

type invoker struct {
	Injector inject.Injector
	Handlers map[reflect.Type]interface{}
	Chains   map[reflect.Type][]*paramGetter
}

func (inv *invoker) makeChain(handler interface{}, injector inject.Injector, deep int) ([]*paramGetter, error) {

	if deep > 10 {
		return nil, errors.New("too deep chain")
	}

	var chain []*paramGetter

	typo := reflect.TypeOf(handler)

	for i := 0; i < typo.NumIn(); i++ {

		in := typo.In(i)

		if injector.Get(in).IsValid() {
			continue
		}

		if inv.Handlers[in] == nil {
			return nil, fmt.Errorf("unknown argument: %v", in)
		}

		inc, err := inv.makeChain(inv.Handlers[in], injector, deep+1)
		if err != nil {
			return nil, err
		}

		chain = append(chain, inc...)
	}

	return append(chain, &paramGetter{Type: typo.Out(0), Handler: handler}), nil
}

func (inv *invoker) makeChains(params ...interface{}) error {

	injector := inject.New()
	injector.SetParent(inv.Injector)

	for _, param := range params {

		typo := reflect.TypeOf(param)
		if typo == nil {
			return errors.New("nil param")
		}

		if inv.Injector.Get(typo).IsValid() || inv.Handlers[typo] != nil {
			return fmt.Errorf("conflict type: %v", typo)
		}

		injector.Map(param)
	}

	for typo, handler := range inv.Handlers {
		if chain, err := inv.makeChain(handler, injector, 1); err != nil {
			return err
		} else {
			inv.Chains[typo] = chain
		}
	}

	return nil
}
