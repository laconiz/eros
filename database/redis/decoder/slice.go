// 切片类型的反序列化

package decoder

import (
	"reflect"
)

func decodeSlice(recv interface{}, replies []interface{}) error {

	// 初始化切片
	typo := reflect.TypeOf(recv).Elem()
	slice := reflect.ValueOf(recv).Elem()
	slice.Set(reflect.MakeSlice(typo, 0, len(replies)))

	for _, reply := range replies {

		// 构造切片数据
		value, err := decodeValue(typo.Elem(), reply)
		if err != nil {
			return err
		}

		// 写入数据
		slice.Set(reflect.Append(slice, value))
	}

	return nil
}
