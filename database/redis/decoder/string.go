// 字符串/字节流/布尔值的反序列化

package decoder

import "github.com/gomodule/redigo/redis"

func decodeString(recv interface{}, reply interface{}) error {

	s, err := redis.String(reply, nil)
	if err != nil {
		return err
	}

	*(recv.(*string)) = s
	return nil
}

func decodeBytes(recv interface{}, reply interface{}) error {

	b, err := redis.Bytes(reply, nil)
	if err != nil {
		return err
	}

	*(recv.(*[]byte)) = b
	return nil
}

func decodeBool(recv interface{}, reply interface{}) error {

	b, err := redis.Bool(reply, nil)
	if err != nil {
		return err
	}

	*(recv.(*bool)) = b
	return nil
}
