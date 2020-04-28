package proto

import "github.com/laconiz/eros/network/message"

// ---------------------------------------------------------------------------------------------------------------------
// 用户ID

type UserID uint64

// ---------------------------------------------------------------------------------------------------------------------
// 性别

type Gender uint8

const (
	GenderInvalid Gender = iota
	GenderMale
	GenderFemale
)

// ---------------------------------------------------------------------------------------------------------------------
// 用户定义

type User struct {
	ID     UserID // ID
	Name   string // 昵称
	Avatar string // 头像
	Gender Gender // 性别
}

// ---------------------------------------------------------------------------------------------------------------------
// 用户信息消息

type ProfileACK struct {
	User  User
	Phone string
}

// ---------------------------------------------------------------------------------------------------------------------
// 修改用户信息消息

type ProfileModifyREQ struct {
	Name   string
	Phone  string
	Avatar string
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	message.Register(ProfileACK{}, message.JsonCodec())
	message.Register(ProfileModifyREQ{}, message.JsonCodec())
}
