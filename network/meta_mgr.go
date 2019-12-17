package network

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
)

func RegisterMeta(msg interface{}, codec Codec) error {

	typo := reflect.TypeOf(msg)

	if typo == nil {
		return errors.New("register a nil message")
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	hash := fnv.New32()
	hash.Write([]byte(typo.String()))

	return RegisterMetaEx(MessageID(hash.Sum32()), msg, codec)
}

func RegisterMetaEx(id MessageID, msg interface{}, codec Codec) error {

	typo := reflect.TypeOf(msg)

	if typo == nil {
		return errors.New("register a nil message")
	}
	if codec == nil {
		return errors.New("register a nil codec")
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	meta := &Meta{
		id:    id,
		typo:  typo,
		codec: codec,
	}

	if metaByID[id] != nil {
		return fmt.Errorf("conflict meta: %v - %v", meta, metaByID[id])
	}
	if metaByType[typo] != nil {
		return fmt.Errorf("conflict meta: %v - %v", meta, metaByType[typo])
	}

	metaByID[id] = meta
	metaByName[typo.String()] = meta
	metaByType[typo] = meta

	logger.Infof("meta: %-12d %-30s %#v", id, typo, codec)
	return nil
}

func MetaByID(id MessageID) *Meta {
	return metaByID[id]
}

func MetaByName(name string) *Meta {
	return metaByName[name]
}

func MetaByType(typo reflect.Type) *Meta {

	if typo == nil {
		return nil
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	return metaByType[typo]
}

func MetaByMsg(msg interface{}) *Meta {
	return MetaByType(reflect.TypeOf(msg))
}

var (
	metaByID   = map[MessageID]*Meta{}
	metaByName = map[string]*Meta{}
	metaByType = map[reflect.Type]*Meta{}
)
