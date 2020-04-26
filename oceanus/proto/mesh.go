package proto

import "github.com/laconiz/eros/network/message"

type MeshID string

type Mesh struct {
	ID    MeshID `json:"id"`              // ID
	Addr  string `json:"addr"`            // 侦听地址
	Proxy string `json:"proxy,omitempty"` // 代理地址
}

type MeshJoin struct {
	Mesh *Mesh `json:"mesh"`
}

type MeshQuit struct {
	ID MeshID `json:"id"`
}

func init() {
	message.Register(MeshJoin{}, message.JsonCodec())
	message.Register(MeshQuit{}, message.JsonCodec())
}
