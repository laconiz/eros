package proto

import "github.com/laconiz/eros/network/message"

// ---------------------------------------------------------------------------------------------------------------------
// 用户ID

type UserID uint64

// ---------------------------------------------------------------------------------------------------------------------
// 性别

type Gender uint8

const (
	Male Gender = iota
	Female
)

// ---------------------------------------------------------------------------------------------------------------------
// 用户定义

type User struct {
	UserID UserID // ID
	Name   string // 昵称
	Avatar string // 头像
	Gender Gender // 性别
}

// ---------------------------------------------------------------------------------------------------------------------
// 用户信息消息

type UserInfoACK struct {
	User
}

// ---------------------------------------------------------------------------------------------------------------------
// 修改用户信息消息

type UserNameChangeREQ struct {
	Name string
}

type UserPhoneChangeREQ struct {
	Phone string
}

type UserAvatarChangeREQ struct {
	Avatar string
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	message.Register(UserInfoACK{}, message.JsonCodec())
	message.Register(UserNameChangeREQ{}, message.JsonCodec())
	message.Register(UserPhoneChangeREQ{}, message.JsonCodec())
	message.Register(UserAvatarChangeREQ{}, message.JsonCodec())
}
