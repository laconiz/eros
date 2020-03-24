package proto

import "github.com/laconiz/eros/network/message"

// ---------------------------------------------------------------------------------------------------------------------
// 物品ID

type ItemID uint32

const (
	ItemCoin   ItemID = 1 // 金币
	ItemTicket ItemID = 2 // 券
)

// ---------------------------------------------------------------------------------------------------------------------
// 物品定义

type Item struct {
	ID  ItemID `json:"id"`  // 物品ID
	Num int64  `json:"num"` // 物品数量
}

// ---------------------------------------------------------------------------------------------------------------------
// 物品列表

type ItemMap map[ItemID]int64

func (m ItemMap) Slice() ItemSlice {
	s := ItemSlice{}
	for id, num := range m {
		s = append(s, &Item{ID: id, Num: num})
	}
	return s
}

type ItemSlice []*Item

func (s ItemSlice) Map() ItemMap {
	m := ItemMap{}
	for _, item := range s {
		m[item.ID] = item.Num
	}
	return m
}

// ---------------------------------------------------------------------------------------------------------------------
// 物品变化类型

type ItemChangeReason uint32

const (
	ItemChangeByAdmin ItemChangeReason = 1001 // 后台修改
	ItemChangeByMail  ItemChangeReason = 1002 // 邮件修改
)

// ---------------------------------------------------------------------------------------------------------------------
// 物品列表消息

type ItemListREQ struct {
}

type ItemListACK struct {
	Items ItemSlice `json:"items,omitempty"`
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	message.Register(ItemListREQ{}, message.JsonCodec())
	message.Register(ItemListACK{}, message.JsonCodec())
}
