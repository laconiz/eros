package proto

import "github.com/laconiz/eros/network/message"

type KickReason int

const (
	KickByServerClosed KickReason = iota // 服务器关闭
	KickByInvalidToken                   // 非法TOKEN
	KickByExpiredToken                   // TOKEN已过期
	KickByOtherLogin                     // 其他地方登录
)

type KickACK struct {
	Reason KickReason `json:"reason"`
}

func init() {
	message.Register(KickACK{}, message.JsonCodec())
}
