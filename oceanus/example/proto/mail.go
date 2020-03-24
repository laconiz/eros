package proto

import "github.com/laconiz/eros/network/message"

// ---------------------------------------------------------------------------------------------------------------------
// 邮件ID

type MailID uint64

// ---------------------------------------------------------------------------------------------------------------------
// 邮件状态

type MailState uint8

const (
	MailStateUnread   MailState = iota // 未读
	MailStateRead                      // 已读
	MailStateReceived                  // 已领取
)

// ---------------------------------------------------------------------------------------------------------------------
// 邮件操作类型

type MailOperation uint8

const (
	MailOperationRead    MailOperation = iota // 读取
	MailOperationReceive                      // 领取
	MailOperationDelete                       // 删除
)

// ---------------------------------------------------------------------------------------------------------------------
// 邮件定义

type Mail struct {
	ID      MailID    `json:"id"`              // ID
	Title   string    `json:"title"`           // 标题
	Content string    `json:"content"`         // 内容
	State   MailState `json:"state"`           // 状态
	Expired int64     `json:"expired"`         // 过期时间
	Items   ItemSlice `json:"items,omitempty"` // 附件
}

// ---------------------------------------------------------------------------------------------------------------------
// 邮件列表消息

type MailListREQ struct {
}

type MailListACK struct {
	Mails []*Mail `json:"mails"` // 全部邮件
}

// ---------------------------------------------------------------------------------------------------------------------
// 邮件操作消息

type MailOperateREQ struct {
	ID        MailID        `json:"id"`        // ID 当ID为0时表示操作所有当前可执行指定操作的邮件
	Operation MailOperation `json:"operation"` // 操作
}

// ---------------------------------------------------------------------------------------------------------------------
// 邮件更新消息

type MailUpdateACK struct {
	Updated []*Mail  `json:"updated,omitempty"` // 更新的邮件(新增/修改)
	Deleted []MailID `json:"deleted,omitempty"` // 删除的邮件ID列表
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	message.Register(MailListREQ{}, message.JsonCodec())
	message.Register(MailListACK{}, message.JsonCodec())
	message.Register(MailOperateREQ{}, message.JsonCodec())
	message.Register(MailUpdateACK{}, message.JsonCodec())
}
