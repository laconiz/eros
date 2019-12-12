// 指定类型的反序列化

package decoder

import "reflect"

func decodeValue(typo reflect.Type, reply interface{}) (reflect.Value, error) {

	// 处理指针类型结构
	ptr := false
	if typo.Kind() == reflect.Ptr {
		ptr = true
		typo = typo.Elem()
	}

	// 构造数据
	var err error
	value := reflect.New(typo)

	switch reply.(type) {
	case nil:
		// 返回值为空时 使用类型的默认值填充
	default:
		// 反序列化数据
		err = Decode(value.Interface(), reply, nil)
	}

	// 非指针类型返回值类型数据
	if !ptr {
		value = value.Elem()
	}

	return value, err
}
