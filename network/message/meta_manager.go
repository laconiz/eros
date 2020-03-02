package message

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
)

type MetaMgr interface {
	Register(msg interface{}, codec Codec) (Meta, error)
	RegisterEx(id ID, name string, msg interface{}, codec Codec) (Meta, error)
	MetaByID(id ID) (Meta, bool)
	MetaByName(name string) (Meta, bool)
	MetaByType(typo reflect.Type) (Meta, bool)
	MetaByMessage(msg interface{}) (Meta, bool)
}

type metaMgr struct {
	idMap   map[ID]Meta
	nameMap map[string]Meta
	typeMap map[reflect.Type]Meta
}

func (m *metaMgr) Register(msg interface{}, codec Codec) (Meta, error) {

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

	return m.RegisterEx(id, name, msg, codec)
}

func (m *metaMgr) RegisterEx(id ID, name string, msg interface{}, codec Codec) (Meta, error) {

	if meta, ok := m.idMap[id]; ok {
		return nil, fmt.Errorf("conflict meta id: %s - %s", name, meta.Name())
	}

	if meta, ok := m.nameMap[name]; ok {
		return nil, fmt.Errorf("conflict meta name: %s - %s", name, meta.Name())
	}

	typo := reflect.TypeOf(msg)
	if typo == nil {
		return nil, fmt.Errorf("register a nil message: %s", name)
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	if meta, ok := m.typeMap[typo]; ok {
		return nil, fmt.Errorf("conflict meta type: %s - %s", name, meta.Name())
	}

	if codec == nil {
		return nil, fmt.Errorf("register meta %s with a nil codec", name)
	}

	meta := &meta{id: id, name: name, typo: typo, codec: codec}

	m.idMap[id] = meta
	m.nameMap[name] = meta
	m.typeMap[typo] = meta

	return meta, nil
}

func (m *metaMgr) MetaByID(id ID) (Meta, bool) {
	meta, ok := m.idMap[id]
	return meta, ok
}

func (m *metaMgr) MetaByName(name string) (Meta, bool) {
	meta, ok := m.nameMap[name]
	return meta, ok
}

func (m *metaMgr) MetaByType(typo reflect.Type) (Meta, bool) {
	meta, ok := m.typeMap[typo]
	return meta, ok
}

func (m *metaMgr) MetaByMessage(msg interface{}) (Meta, bool) {

	typo := reflect.TypeOf(msg)
	if typo == nil {
		return nil, false
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	meta, ok := m.typeMap[typo]
	return meta, ok
}

var globalMetaMgr = NewMetaMgr()

func Register(msg interface{}, codec Codec) (Meta, error) {
	return globalMetaMgr.Register(msg, codec)
}

func RegisterEx(id ID, name string, msg interface{}, codec Codec) (Meta, error) {
	return globalMetaMgr.RegisterEx(id, name, msg, codec)
}

func MetaByID(id ID) (Meta, bool) {
	return globalMetaMgr.MetaByID(id)
}

func MetaByName(name string) (Meta, bool) {
	return globalMetaMgr.MetaByName(name)
}

func MetaByType(typo reflect.Type) (Meta, bool) {
	return globalMetaMgr.MetaByType(typo)
}

func MetaByMessage(msg interface{}) (Meta, bool) {
	return globalMetaMgr.MetaByMessage(msg)
}

func NewMetaMgr() MetaMgr {
	return &metaMgr{
		idMap:   map[ID]Meta{},
		nameMap: map[string]Meta{},
		typeMap: map[reflect.Type]Meta{},
	}
}
