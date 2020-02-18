// 依赖注入

package ioc

import (
	"fmt"
	"reflect"

	"github.com/codegangsta/inject"
)

type Injector = inject.Injector

// 依赖注入接口
type Invoker func(params ...interface{}) (arguments []interface{}, result []reflect.Value, err error)

// ---------------------------------------------------------------------------------------------------------------------

// 创建一个参数解释器
func New() *Squirt {
	return &Squirt{injector: inject.New(), builders: map[Type]*Builder{}}
}

// 注入器
type Squirt struct {
	params   []interface{}     // 基础参数
	injector Injector          // 参数注入器
	builders map[Type]*Builder // 参数生成函数
	err      error             // 注入器错误
}

// 拷贝一个注入器
func (squirt *Squirt) Copy() *Squirt {

	n := New()

	n.err = squirt.err

	for _, param := range squirt.params {
		n.params = append(n.params, param)
		n.injector.Map(param)
	}

	for key, value := range squirt.builders {
		n.builders[key] = value
	}

	return n
}

// 添加基础参数
func (squirt *Squirt) Params(params ...interface{}) *Squirt {

	return squirt.safeExecute(func() error {

		for _, param := range params {

			typo := reflect.TypeOf(param)
			if typo == nil {
				return fmt.Errorf("invalid param: %v", typo)
			}

			if squirt.injector.Get(typo).IsValid() {
				return fmt.Errorf("conflict param: %v", typo)
			}

			squirt.params = append(squirt.params, param)
			squirt.injector.Map(param)
		}
		return nil
	})
}

// 添加参数生成器
func (squirt *Squirt) Creators(creators ...interface{}) *Squirt {

	return squirt.safeExecute(func() error {

		for _, creator := range creators {

			builder, err := newBuilder(creator)
			if err != nil {
				return err
			}

			if squirt.hasArg(builder.typo) {
				return fmt.Errorf("conflict creator: %v", builder.typo)
			}

			squirt.builders[builder.typo] = builder
		}

		return nil
	})
}

// 添加参数生成器
func (squirt *Squirt) Creator(typo Type, creator interface{}) *Squirt {

	return squirt.safeExecute(func() error {

		if squirt.hasArg(typo) {
			return fmt.Errorf("conflict creator: %v", typo)
		}

		builder, err := newBuilder(creator)
		if err != nil {
			return err
		}

		builder.typo = typo
		squirt.builders[typo] = builder
		return nil
	})
}

// 当注入器状态正常时执行函数
func (squirt *Squirt) safeExecute(f func() error) *Squirt {
	if squirt.err == nil {
		squirt.err = f()
	}
	return squirt
}

// 查询参数
func (squirt *Squirt) hasArg(typo reflect.Type) bool {
	return squirt.injector.Get(typo).IsValid() || squirt.builders[typo] != nil
}

// 查询状态
func (squirt *Squirt) Error() error {
	return squirt.err
}

// 生成一个接口
func (squirt *Squirt) Handle(handler interface{}, args ...interface{}) (Invoker, error) {

	if squirt.err != nil {
		return nil, squirt.err
	}

	typo, err := checkHandler(handler)
	if err != nil {
		return nil, err
	}

	injector := inject.New()
	injector.SetParent(squirt.injector)
	for _, arg := range args {
		injector.Map(arg)
	}

	var chain Chain
	for i := 0; i < typo.NumIn(); i++ {

		in := typo.In(i)
		if injector.Get(in).IsValid() {
			continue
		}

		subChain, err := squirt.makeChain(in, injector, map[Type]bool{})
		if err != nil {
			return nil, err
		}
		chain = append(chain, subChain...)
	}

	return squirt.makeInvoker(distinctChain(chain), handler), nil
}

// 生成依赖注入调用接口
func (squirt *Squirt) makeInvoker(chain Chain, handler interface{}) Invoker {

	return func(params ...interface{}) ([]interface{}, []reflect.Value, error) {

		injector := inject.New()
		injector.SetParent(squirt.injector)
		for _, param := range params {
			injector.Map(param)
		}

		var arguments []interface{}

		for _, builder := range chain {

			values, err := injector.Invoke(builder.creator)
			if err != nil {
				return nil, nil, fmt.Errorf("invoke creator[%v] error: %w", builder.typo, err)
			}

			if !values[1].IsNil() {
				err := values[1].Interface().(error)
				return nil, nil, fmt.Errorf("build %v error: %w", builder.typo, err)
			}

			if values[0].CanInterface() {
				injector.Map(values[0].Interface())
				arguments = append(arguments, values[0].Interface())
			}
		}

		values, err := injector.Invoke(handler)
		if err != nil {
			return nil, nil, fmt.Errorf("invoke handler error: %w", err)
		}
		return arguments, values, nil
	}
}

// 生成指定类型的参数生成链
func (squirt *Squirt) makeChain(typo Type, injector Injector, paths map[Type]bool) (Chain, error) {

	builder, ok := squirt.builders[typo]
	if !ok {
		return nil, fmt.Errorf("unknown argument: %v", typo)
	}

	if _, ok := paths[builder.typo]; ok {
		return nil, fmt.Errorf("circle creator: %v", typo)
	}

	var chain Chain

	creatorType := reflect.TypeOf(builder.creator)
	for i := 0; i < creatorType.NumIn(); i++ {

		in := creatorType.In(i)
		if injector.Get(in).IsValid() {
			continue
		}

		newPaths := map[Type]bool{typo: true}
		for key, value := range paths {
			newPaths[key] = value
		}

		subChain, err := squirt.makeChain(in, injector, newPaths)
		if err != nil {
			return nil, err
		}
		chain = append(chain, subChain...)
	}

	return append(chain, builder), nil
}

// 获取依赖注入器中不存在的参数
func (squirt *Squirt) UnknownArgs(handler interface{}, params ...interface{}) ([]Type, error) {

	typo, err := checkHandler(handler)
	if err != nil {
		return nil, err
	}

	injector := inject.New()
	injector.SetParent(squirt.injector)
	for _, param := range params {
		injector.Map(param)
	}

	var args []Type
	for i := 0; i < typo.NumIn(); i++ {
		in := typo.In(i)
		if !injector.Get(in).IsValid() && squirt.builders[in] == nil {
			args = append(args, in)
		}
	}

	return args, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// 获取接口类型
func checkHandler(handler interface{}) (Type, error) {

	typo := reflect.TypeOf(handler)
	if typo == nil || typo.Kind() != reflect.Func {
		return nil, fmt.Errorf("invalid handler type: %v", typo)
	}

	return typo, nil
}

// 参数生成链去重
func distinctChain(chain Chain) Chain {

	var result Chain
	types := map[Type]bool{}

	for _, builder := range chain {
		if !types[builder.typo] {
			types[builder.typo] = true
			result = append(result, builder)
		}
	}

	return result
}
