package model

import (
	"errors"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// 各种DB统一查询为空返回值
// 把redis/sql/elastic各自的空查询错误接口统一

var RecordNotFound = errors.New("record not found")

// ---------------------------------------------------------------------------------------------------------------------
// SQL模型定义

type Model struct {
	ID        uint       `json:"primary";gorm:"primary_key"` // 主键ID
	CreatedAt time.Time  `json:"created"`                    // 创建时间
	UpdatedAt time.Time  `json:"updated"`                    // 更新时间
	DeletedAt *time.Time `json:"deleted";sql:"index"`        // 删除时间
}

// ---------------------------------------------------------------------------------------------------------------------
// redis key定义

const (
	UserMaxID  = "user.maxID" // 当前最大用户ID
	ItemPrefix = "item."      // 物品表前缀
)

// ---------------------------------------------------------------------------------------------------------------------
// elastic 索引&别名定义

const (
	UserLoginLogAlias   = "user_login_log"         // 用户登录日志别名
	UserLoginLogPrefix  = UserLoginLogAlias + "_"  // 用户登录日志索引前缀
	ItemChangeLogAlias  = "item_change_log"        // 物品改变日志别名
	ItemChangeLogPrefix = ItemChangeLogAlias + "_" // 物品改变日志索引前缀
)

// ---------------------------------------------------------------------------------------------------------------------
// model 模块名

const (
	ModuleUser = "user"
	ModuleItem = "item"
	ModuleMail = "mail"
)

// ---------------------------------------------------------------------------------------------------------------------
// model模块日志接口生成器

const (
	nameModel = "model"
	fieldMgr  = "manager"
)

func Logger(name string) logis.Logger {
	return logisor.Field(logis.Module, nameModel).Field(fieldMgr, name)
}
