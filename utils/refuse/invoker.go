// 参数注入器

package refuse

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/codegangsta/inject"
)

func New(params ...interface{}) (*Squirt, error) {

	injector := inject.New()
	handlers := map[reflect.Type]interface{}{}

	for _, param := range params {
		// 检查参数类型
		typo := reflect.TypeOf(param)
		if typo == nil {
			return nil, errors.New("nil param")
		}
		// 写入基础参数
		if typo.Kind() != reflect.Func {
			// 重复参数类型
			if injector.Get(typo).IsValid() {
				return nil, fmt.Errorf("conflict type: %v", typo)
			}
			injector.Map(param)
			continue
		}
		// 检查参数生成器
		if typo.NumOut() != 2 || typo.Out(1).Name() != "error" {
			return nil, fmt.Errorf("invalid func: %v", param)
		}
		// 重复参数类型
		out := typo.Out(0)
		if injector.Get(out).IsValid() || handlers[out] != nil {
			return nil, fmt.Errorf("conflict type: %v", out)
		}
		handlers[out] = param
	}

	return &Squirt{
		injector: injector,
		handlers: handlers,
		chains:   map[reflect.Type]Chain{},
	}, nil
}

// 参数生成器
type Func struct {
	typo    reflect.Type
	handler interface{}
}

type Chain []*Func

// 注入器
type Squirt struct {
	injector inject.Injector
	handlers map[reflect.Type]interface{}
	chains   map[reflect.Type]Chain
}

func (s *Squirt) MakeChains(params ...interface{}) (*Squirt, error) {

	injector := inject.New()
	injector.SetParent(s.injector)

	for _, param := range params {
		// 检查参数类型
		typo := reflect.TypeOf(param)
		if typo == nil {
			return nil, errors.New("nil param")
		}
		// 重复参数类型
		if s.injector.Get(typo).IsValid() || s.handlers[typo] != nil {
			return nil, fmt.Errorf("conflict type: %v", typo)
		}
		injector.Map(param)
	}

}

func (s *Squirt) makeChain(typo reflect.Type, injector inject.Injector, depth int) error {

	if _, ok := s.chains[typo]; ok {
		return nil
	}

	if depth >= maxDepth {
		return errors.New("too deep function")
	}

	handler := s.handlers[typo]

	var chain Chain
	for i := 0; i < typo.NumIn(); i++ {
		in := typo.In(i)
		// 参数可以直接访问
		if injector.Get(in).IsValid() {
			continue
		}
		// 已生成的参数生成链
		if _, ok := s.chains[in]; !ok {
			if err := s.makeChain(in, injector, depth+1); err != nil {
				return err
			}
		}
		chain = append(chain, s.chains[in]...)
	}

}

const (
	defaultDepth = 0
	maxDepth     = 10
)
