package ioc

import (
	"fmt"
	"reflect"
)

type Type = reflect.Type

// ---------------------------------------------------------------------------------------------------------------------

// 生成参数生成器
func newBuilder(creator interface{}) (*Builder, error) {

	typo := reflect.TypeOf(creator)
	if typo == nil || typo.Kind() != reflect.Func {
		return nil, fmt.Errorf("creator must be func, got %v", typo)
	}

	et := reflect.TypeOf((*error)(nil)).Elem()
	if typo.NumOut() != 2 || !typo.Out(1).Implements(et) {
		return nil, fmt.Errorf("creator must be func(...any) (any, error), got %v", typo)
	}

	return &Builder{typo: typo.Out(0), creator: creator}, nil
}

// 参数生成器
type Builder struct {
	typo    Type        // 生成的参数类型
	creator interface{} // 参数生成函数
}

// ---------------------------------------------------------------------------------------------------------------------

// 参数生成链
type Chain []*Builder
