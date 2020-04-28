package model

import (
	"github.com/laconiz/eros/oceanus/example/proto"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

type ItemID = proto.ItemID
type ItemChangeReason = proto.ItemChangeReason

const (
	ItemCoin = proto.ItemCoin
)

// ---------------------------------------------------------------------------------------------------------------------
// 物品变化日志 ELASTIC

type ItemChangeLog struct {
	Item    ItemID           `json:"item"`    // 物品ID
	User    UserID           `json:"user"`    // 用户ID
	Value   int64            `json:"value"`   // 变化数量
	Latest  int64            `json:"latest"`  // 最新数量
	Reason  ItemChangeReason `json:"reason"`  // 变化原因
	Success bool             `json:"success"` // 是否成功
	Time    time.Time        `json:"time"`    // 变化时间
}
