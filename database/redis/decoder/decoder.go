// redigo返回数据的反序列化

package decoder

import (
	"fmt"
	"reflect"

	"github.com/gomodule/redigo/redis"

	"github.com/laconiz/eros/utils/json"
)

// 反序列化redigo的返回数据
// 接收参数必须为指针
func Decode(recv interface{}, reply interface{}, err error) error {

	// 检测错误
	if err != nil {
		return err
	}

	// 序列化简单类型
	if ok, err := decodeSimple(recv, reply); ok {
		return err
	}

	// 检查返回值参数类型
	typo := reflect.TypeOf(recv)
	if typo == nil || typo.Kind() != reflect.Ptr {
		return fmt.Errorf("receiver must be pointer, got %v", typo)
	}

	// 如果返回数据为字节流 直接使用json进行反序列化
	if bytes, err := redis.Bytes(reply, nil); err == nil {
		return json.Unmarshal(bytes, recv)
	}

	// 获取切片数据
	replies, err := redis.Values(reply, nil)
	if err != nil {
		return err
	}

	typo = typo.Elem()

	// 反序列化复杂类型
	switch typo.Kind() {

	case reflect.Slice:
		return decodeSlice(recv, replies)

	case reflect.Map:
		return decodeMap(recv, replies)

	case reflect.Struct:
		return decodeStruct(recv, replies)
	}

	return fmt.Errorf("unsupported type: %v", typo)
}

// 序列化内置类型
func decodeSimple(recv interface{}, reply interface{}) (bool, error) {

	switch recv.(type) {
	case *string:
		return true, decodeString(recv, reply)
	case *[]byte:
		return true, decodeBytes(recv, reply)
	case *bool:
		return true, decodeBool(recv, reply)
	case *int8:
		return true, decodeInt8(recv, reply)
	case *int16:
		return true, decodeInt16(recv, reply)
	case *int32:
		return true, decodeInt32(recv, reply)
	case *int64:
		return true, decodeInt64(recv, reply)
	case *int:
		return true, decodeInt(recv, reply)
	case *uint8:
		return true, decodeUint8(recv, reply)
	case *uint16:
		return true, decodeUint16(recv, reply)
	case *uint32:
		return true, decodeUint32(recv, reply)
	case *uint64:
		return true, decodeUint64(recv, reply)
	case *uint:
		return true, decodeUint(recv, reply)
	}

	return false, nil
}
