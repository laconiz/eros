// 依赖注入

package ioc

import (
	"reflect"

	"github.com/codegangsta/inject"
)

// 创建一个参数解释器
func New() *Squirt {
	return &Squirt{
		injector: inject.New(),
		builders: map[reflect.Type]*Builder{},
	}
}

// 依赖注入接口
type Invoker func(...interface{}) ([]interface{}, []reflect.Value, error)

// 注入器
type Squirt struct {
	// 基础参数
	params []interface{}
	//
	injector inject.Injector
	// 参数生成函数
	builders map[reflect.Type]*Builder
	//
	err error
}

// 拷贝一个注入器
func (s *Squirt) Copy() *Squirt {
	n := New()
	n.err = s.err
	if n.err == nil {
		for _, param := range s.params {
			n.params = append(n.params, param)
			n.injector.Map(param)
		}
		for key, value := range s.builders {
			n.builders[key] = value
		}
	}
	return n
}

// 添加基础参数
func (s *Squirt) Params(params ...interface{}) *Squirt {
	return s.safeExecute(func() error {
		for _, param := range params {
			// 检查参数类型
			typo := reflect.TypeOf(param)
			if typo == nil {
				return &Error{msg: "nil param", typo: typo}
			}
			// 重复参数类型
			if s.injector.Get(typo).IsValid() {
				return &Error{msg: "conflict param", typo: typo}
			}
			s.params = append(s.params, param)
			s.injector.Map(param)
		}
		return nil
	})
}

// 添加参数生成器
func (s *Squirt) Functions(functions ...interface{}) *Squirt {
	return s.safeExecute(func() error {
		for _, function := range functions {
			// 构造参数生成器
			builder, err := newBuilder(function)
			if err != nil {
				return err
			}
			// 检查重复参数
			if s.hasArg(builder.typo) {
				return &Error{msg: "conflict builder", typo: builder.typo}
			}
			s.builders[builder.typo] = builder
		}
		return nil
	})
}

// 添加参数生成器
func (s *Squirt) Function(typo reflect.Type, function interface{}) *Squirt {
	return s.safeExecute(func() error {
		// 检查重复参数
		if s.hasArg(typo) {
			return &Error{msg: "conflict builder", typo: typo}
		}
		// 构造参数生成器
		builder, err := newBuilder(function)
		if err != nil {
			return err
		}
		// 重置参数类型
		builder.typo = typo
		s.builders[typo] = builder
		return nil
	})
}

// 当注入器状态正常时执行函数
func (s *Squirt) safeExecute(f func() error) *Squirt {
	if s.err == nil {
		s.err = f()
	}
	return s
}

// 查询参数
func (s *Squirt) hasArg(typo reflect.Type) bool {
	return s.injector.Get(typo).IsValid() || s.builders[typo] != nil
}

// 查询状态
func (s *Squirt) Error() error {
	return s.err
}

func (s *Squirt) Handle(handler interface{}, args ...interface{}) (Invoker, error) {
	// 校验接口类型
	typo, err := checkHandler(handler)
	if err != nil {
		return nil, err
	}
	// 基础参数
	injector := inject.New()
	injector.SetParent(s.injector)
	for _, arg := range args {
		injector.Map(arg)
	}
	// 参数生成链
	chain := Chain{}
	for i := 0; i < typo.NumIn(); i++ {
		in := typo.In(i)
		if injector.Get(in).IsValid() {
			continue
		}
		// 构造参数生成链
		c, err := s.makeChain(in, injector, map[reflect.Type]bool{})
		if err != nil {
			return nil, err
		}
		chain = append(chain, c...)
	}
	// 构造调用接口
	return s.invoker(chain.Distinct(), handler), nil
}

//
func (s *Squirt) invoker(chain Chain, handler interface{}) Invoker {

	return func(args ...interface{}) ([]interface{}, []reflect.Value, error) {

		// 构造基础参数
		injector := inject.New()
		injector.SetParent(s.injector)
		for _, arg := range args {
			injector.Map(arg)
		}

		var arguments []interface{}

		// 调用参数生成链
		for _, builder := range chain {

			values, err := injector.Invoke(builder.function)
			if err != nil {
				return nil, nil, err
			}

			if !values[1].IsNil() {
				return nil, nil, values[1].Interface().(error)
			}

			if values[0].CanInterface() {
				injector.Map(values[0].Interface())
				arguments = append(arguments, values[0].Interface())
			}
		}

		// 调用接口
		values, err := injector.Invoke(handler)
		return arguments, values, err
	}
}

// 生成指定类型的参数生成链
func (s *Squirt) makeChain(typo reflect.Type, injector inject.Injector, paths map[reflect.Type]bool) (Chain, error) {
	// 获取参数生成器
	builder, ok := s.builders[typo]
	if !ok {
		return nil, &Error{msg: "unknown argument", typo: typo}
	}
	// 循环链
	if _, ok := paths[builder.typo]; ok {
		return nil, &Error{msg: "circle builder", typo: builder.typo}
	}
	// 检索参数列表
	var chain Chain
	functionType := reflect.TypeOf(builder.function)
	for i := 0; i < functionType.NumIn(); i++ {
		in := functionType.In(i)
		// 参数可以直接访问
		if injector.Get(in).IsValid() {
			continue
		}
		// 构造新路径
		ps := map[reflect.Type]bool{builder.typo: true}
		for key, value := range paths {
			ps[key] = value
		}
		// 构造子链
		c, err := s.makeChain(in, injector, ps)
		if err != nil {
			return nil, err
		}
		chain = append(chain, c...)
	}
	return append(chain, builder), nil
}

// 获取参数生成器中没有的参数
func (s *Squirt) UnknownArgs(handler interface{}) ([]reflect.Type, error) {
	// 校验接口类型
	typo, err := checkHandler(handler)
	if err != nil {
		return nil, err
	}
	// 未知参数列表
	var args []reflect.Type
	// 获取未知参数
	for i := 0; i < typo.NumIn(); i++ {
		in := typo.In(i)
		if !s.hasArg(in) {
			args = append(args, in)
		}
	}
	return args, nil
}

// 获取接口类型
func checkHandler(handler interface{}) (reflect.Type, error) {
	typo := reflect.TypeOf(handler)
	if typo == nil || typo.Kind() != reflect.Func {
		return nil, &Error{msg: "invalid handler", typo: typo}
	}
	return typo, nil
}
