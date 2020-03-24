package model

import (
	"github.com/laconiz/eros/oceanus/example/proto"
)

// ---------------------------------------------------------------------------------------------------------------------

type UserID = proto.UserID
type Gender = proto.Gender

// ---------------------------------------------------------------------------------------------------------------------
// 用户定义 SQL

type User struct {
	Model
	UserID   UserID `json:"id";gorm:"unique_index"`    // ID
	Name     string `json:"name";gorm:"unique_index"`  // 昵称
	Avatar   string `json:"avatar"`                    // 头像
	Gender   Gender `json:"gender"`                    // 性别
	Phone    string `json:"phone";gorm:"unique_index"` // 手机
	Password string `json:"password"`                  // 密码
}
