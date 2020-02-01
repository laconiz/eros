// map类型的反序列化

package decoder

import (
	"fmt"
	"reflect"
)

func decodeMap(recv interface{}, replies []interface{}) error {

	// 检查返回数据数量
	if len(replies)%2 != 0 {
		return fmt.Errorf("replies's length must be dual, got %d", len(replies))
	}

	// 初始化map
	typo := reflect.TypeOf(recv).Elem()
	mapValue := reflect.ValueOf(recv).Elem()
	mapValue.Set(reflect.MakeMap(typo))

	for i := 0; i < len(replies); i += 2 {

		// 构造key
		index, err := decodeValue(typo.Key(), replies[i])
		if err != nil {
			return err
		}

		// 构造value
		value, err := decodeValue(typo.Elem(), replies[i+1])
		if err != nil {
			return err
		}

		// 设置值
		mapValue.SetMapIndex(index, value)
	}

	return nil
}
