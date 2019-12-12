// 结构体的反序列化

package decoder

import (
	"fmt"
	"reflect"
)

func decodeStruct(recv interface{}, replies []interface{}) error {

	typo := reflect.TypeOf(recv).Elem()

	// 检查返回数据数量
	if len(replies) != typo.NumField() {
		return fmt.Errorf("replies's length does not equal struct's field num: %d != %d",
			len(replies), typo.NumField())
	}

	structValue := reflect.ValueOf(recv).Elem()

	for i := 0; i < typo.NumField(); i++ {

		// 构造字段值
		value, err := decodeValue(typo.Field(i).Type, replies[i])
		if err != nil {
			return err
		}

		// 写入数据
		structValue.Field(i).Set(value)
	}

	return nil
}
