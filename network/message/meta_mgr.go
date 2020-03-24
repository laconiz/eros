// meta manager

package message

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
)

// ---------------------------------------------------------------------------------------------------------------------

func Json(msg interface{}) (Meta, error) {
	return Register(msg, JsonCodec())
}

// ---------------------------------------------------------------------------------------------------------------------

func Register(msg interface{}, codec Codec) (Meta, error) {

	typo := reflect.TypeOf(msg)
	if typo == nil {
		return nil, errors.New("register a nil message")
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	name := typo.String()

	hash := fnv.New32()
	hash.Write([]byte(name))
	id := ID(hash.Sum32())

	return RegisterEx(id, name, msg, codec)
}

func RegisterEx(id ID, name string, msg interface{}, codec Codec) (Meta, error) {

	if meta, ok := metaByID[id]; ok {
		return nil, fmt.Errorf("conflict meta id: %s - %s", name, meta.Name())
	}

	if meta, ok := metaByName[name]; ok {
		return nil, fmt.Errorf("conflict meta name: %s - %s", name, meta.Name())
	}

	typo := reflect.TypeOf(msg)
	if typo == nil {
		return nil, fmt.Errorf("register a nil message: %s", name)
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	if meta, ok := metaByType[typo]; ok {
		return nil, fmt.Errorf("conflict meta type: %s - %s", name, meta.Name())
	}

	if codec == nil {
		return nil, fmt.Errorf("register meta %s with a nil codec", name)
	}

	meta := &meta{id: id, name: name, typo: typo, codec: codec}

	metaByID[id] = meta
	metaByName[name] = meta
	metaByType[typo] = meta

	return meta, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func MetaByID(id ID) (Meta, bool) {
	meta, ok := metaByID[id]
	return meta, ok
}

func MetaByName(name string) (Meta, bool) {
	meta, ok := metaByName[name]
	return meta, ok
}

func MetaByType(typo reflect.Type) (Meta, bool) {
	meta, ok := metaByType[typo]
	return meta, ok
}

func MetaByMsg(msg interface{}) (Meta, bool) {

	typo := reflect.TypeOf(msg)
	if typo == nil {
		return nil, false
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	meta, ok := metaByType[typo]
	return meta, ok
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	metaByID   = map[ID]Meta{}
	metaByName = map[string]Meta{}
	metaByType = map[reflect.Type]Meta{}
)
