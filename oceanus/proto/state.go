package proto

import "github.com/laconiz/eros/network/message"

type State struct {
	Version uint32 `json:"ver"`
	Power   int64  `json:"power"`
	Limit   int64  `json:"limit"`
}

func (state *State) Load() int64 {
	return state.Limit - state.Power
}

func init() {
	message.Register(State{}, message.JsonCodec())
}
