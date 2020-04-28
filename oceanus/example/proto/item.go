package proto

import "github.com/laconiz/eros/network/message"

// ---------------------------------------------------------------------------------------------------------------------
// 物品ID

type ItemID string

const (
	ItemCoin   ItemID = "coin"   // 金币
	ItemTicket ItemID = "ticket" // 券
)

// ---------------------------------------------------------------------------------------------------------------------
// 物品定义

type Item struct {
	ID  ItemID // ID
	Num int64  // 数量
}

// ---------------------------------------------------------------------------------------------------------------------
// 物品列表

type ItemMap map[ItemID]int64

func (hash ItemMap) Slice() ItemSlice {

	slice := ItemSlice{}

	for id, num := range hash {
		item := &Item{ID: id, Num: num}
		slice = append(slice, item)
	}

	return slice
}

type ItemSlice []*Item

func (slice ItemSlice) Map() ItemMap {

	hash := ItemMap{}

	for _, item := range slice {
		hash[item.ID] = item.Num
	}

	return hash
}

// ---------------------------------------------------------------------------------------------------------------------
// 物品变化类型

type ItemChangeReason string

const (
	ItemChangeByAdmin ItemChangeReason = "admin" // 后台修改
	ItemChangeByMail  ItemChangeReason = "mail"  // 邮件修改
)

// ---------------------------------------------------------------------------------------------------------------------
// 物品列表消息

type ItemsACK struct {
	Items ItemMap `json:"items,omitempty"`
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	message.Register(ItemsACK{}, message.JsonCodec())
}
