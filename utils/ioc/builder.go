package ioc

import (
	"fmt"
	"reflect"
)

type Error struct {
	msg  string
	typo reflect.Type
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %v", e.msg, e.typo)
}

// 生成参数生成器
func newBuilder(function interface{}) (*Builder, error) {
	// 检查参数类型
	typo := reflect.TypeOf(function)
	if typo == nil || typo.Kind() != reflect.Func {
		return nil, &Error{msg: "need func(...any) (any, error), got", typo: typo}
	}
	// 检查返回值类型
	et := reflect.TypeOf((*error)(nil)).Elem()
	if typo.NumOut() != 2 || !typo.Out(1).Implements(et) {
		return nil, &Error{msg: "need func(...any) (any, error), got", typo: typo}
	}
	return &Builder{typo: typo.Out(0), function: function}, nil
}

type Builder struct {
	typo     reflect.Type
	function interface{}
}

type Chain []*Builder

// 参数生成链去重
func (c Chain) Distinct() Chain {
	if c == nil {
		return c
	}
	var chain Chain
	types := map[reflect.Type]bool{}
	// 去重
	for _, builder := range c {
		if !types[builder.typo] {
			chain = append(chain, builder)
			types[builder.typo] = true
		}
	}
	return chain
}
