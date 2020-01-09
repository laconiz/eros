package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/json"
)

func KV() *kv {
	return &kv{KV: client.KV()}
}

type kv struct {
	*api.KV
}

func (kv *kv) Load(key string, value interface{}) error {

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(pair.Value, value)
}

func (kv *kv) Store(key string, value interface{}) error {

	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	pair := &api.KVPair{
		Key:   key,
		Value: raw,
	}

	_, err = kv.Put(pair, nil)
	return err
}

func (kv *kv) Delete(key string) error {
	_, err := kv.KV.Delete(key, nil)
	return err
}
