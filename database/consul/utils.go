package consul

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/consul/api"

	"github.com/laconiz/eros/utils/json"
)

// 将键值对列表序列化成字典
// 字典的key为键值对的key去掉前缀后的值
// 当strict参数为true时, 序列化失败键值对值时会返回错误
func ParsePairs(prefix string, pairs api.KVPairs, receiver interface{}, strict bool) error {

	typo := reflect.TypeOf(receiver)
	if typo == nil || typo.Kind() != reflect.Ptr {
		return fmt.Errorf("receiver must be pointer, got %v", typo)
	}
	typo = typo.Elem()
	if typo.Kind() != reflect.Map {
		return errors.New("receiver must be a map pointer")
	}

	mapValue := reflect.ValueOf(receiver).Elem()
	mapValue.Set(reflect.MakeMap(typo))

	if typo.Key().Kind() != reflect.String {
		return errors.New("receiver's key must be string")
	}

	for _, pair := range pairs {

		key := strings.Replace(pair.Key, prefix, "", 1)

		value := reflect.New(typo.Elem())
		if err := json.Unmarshal(pair.Value, value.Interface()); err != nil {
			if strict {
				return err
			}
			continue
		}

		mapValue.SetMapIndex(reflect.ValueOf(key), value.Elem())
	}

	return nil
}
