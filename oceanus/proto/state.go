package proto

import "github.com/laconiz/eros/network/message"

type Version int64

type State struct {
	Version Version `json:"version"`
	Credit  int64   `json:"credit"`
}

func init() {
	message.Register(State{}, message.JsonCodec())
}
