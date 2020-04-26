package chat

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/proto"
)

// ---------------------------------------------------------------------------------------------------------------------

type Room struct {
	id      RoomID                  // 房间ID
	users   map[UserID]proto.NodeID // 用户列表
	process oceanus.Oceanus
}

// ---------------------------------------------------------------------------------------------------------------------

func (room *Room) Proto() *RoomACK {

	var list []UserID
	for id := range room.users {
		list = append(list, id)
	}

	return &RoomACK{ID: room.id, Users: list}
}

// ---------------------------------------------------------------------------------------------------------------------
// 向房间内的所有用户发送消息

func (room *Room) Broadcast(msg interface{}) {

	var list []proto.NodeID
	for _, id := range room.users {
		list = append(list, id)
	}

	room.process.SendByIDs(list, msg)
}
