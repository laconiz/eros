// 数值类型的反序列化

package decoder

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func decodeInt8(recv interface{}, reply interface{}) error {

	i, err := redis.Int64(reply, nil)
	if err != nil {
		return err
	}

	value := int8(i)
	if int64(value) != i {
		return fmt.Errorf("decode %d to int8", i)
	}

	*(recv.(*int8)) = value
	return nil
}

func decodeInt16(recv interface{}, reply interface{}) error {

	i, err := redis.Int64(reply, nil)
	if err != nil {
		return err
	}

	value := int16(i)
	if int64(value) != i {
		return fmt.Errorf("decode %d to int16", i)
	}

	*(recv.(*int16)) = value
	return nil
}

func decodeInt32(recv interface{}, reply interface{}) error {

	i, err := redis.Int64(reply, nil)
	if err != nil {
		return err
	}

	value := int32(i)
	if int64(value) != i {
		return fmt.Errorf("decode %d to int32", i)
	}

	*(recv.(*int32)) = value
	return nil
}

func decodeInt64(recv interface{}, reply interface{}) error {

	value, err := redis.Int64(reply, nil)
	if err != nil {
		return err
	}

	*(recv.(*int64)) = value
	return nil
}

func decodeInt(recv interface{}, reply interface{}) error {

	i, err := redis.Int64(reply, nil)
	if err != nil {
		return err
	}

	value := int(i)
	if int64(value) != i {
		return fmt.Errorf("decode %d to int", i)
	}

	*(recv.(*int)) = value
	return nil
}

func decodeUint8(recv interface{}, reply interface{}) error {

	i, err := redis.Uint64(reply, nil)
	if err != nil {
		return err
	}

	value := uint8(i)
	if uint64(value) != i {
		return fmt.Errorf("decode %d to int8", i)
	}

	*(recv.(*uint8)) = value
	return nil
}

func decodeUint16(recv interface{}, reply interface{}) error {

	i, err := redis.Uint64(reply, nil)
	if err != nil {
		return err
	}

	value := uint16(i)
	if uint64(value) != i {
		return fmt.Errorf("decode %d to uint16", i)
	}

	*(recv.(*uint16)) = value
	return nil
}

func decodeUint32(recv interface{}, reply interface{}) error {

	i, err := redis.Uint64(reply, nil)
	if err != nil {
		return err
	}

	value := uint32(i)
	if uint64(value) != i {
		return fmt.Errorf("decode %d to uint32", i)
	}

	*(recv.(*uint32)) = value
	return nil
}

func decodeUint64(recv interface{}, reply interface{}) error {

	value, err := redis.Uint64(reply, nil)
	if err != nil {
		return err
	}

	*(recv.(*uint64)) = value
	return nil
}

func decodeUint(recv interface{}, reply interface{}) error {

	i, err := redis.Uint64(reply, nil)
	if err != nil {
		return err
	}

	value := uint(i)
	if uint64(value) != i {
		return fmt.Errorf("decode %d to uint", i)
	}

	*(recv.(*uint)) = value
	return nil
}
