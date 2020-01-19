// consul键值对操作

package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/json"
	"reflect"
	"strings"
)

var ErrNotFound = errors.New("key not found")

func KV() *kv {
	return &kv{KV: client.KV()}
}

type kv struct {
	*api.KV
}

// 获取键值对
func (kv *kv) Load(key string, value interface{}) error {
	// 加载数据
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return err
	}
	// 没有找到数据
	if pair == nil {
		return ErrNotFound
	}
	// 反序列化数据
	return json.Unmarshal(pair.Value, value)
}

// 存储键值对
func (kv *kv) Store(key string, value interface{}) error {
	// 序列化数据
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// 写入数据
	pair := &api.KVPair{Key: key, Value: raw}
	_, err = kv.Put(pair, nil)
	return err
}

// 删除键值对
func (kv *kv) Delete(key string) error {
	_, err := kv.KV.Delete(key, nil)
	return err
}

// 获取前缀下的所有键值对
func (kv *kv) Loads(prefix string, receiver interface{}, strict bool) error {
	// 获取列表
	pairs, _, err := kv.KV.List(prefix, nil)
	if err != nil {
		return err
	}
	// 反序列化数据
	return ParsePairs(prefix, pairs, receiver, strict)
}
