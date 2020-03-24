package chat

import "github.com/laconiz/eros/oceanus"

// ---------------------------------------------------------------------------------------------------------------------

type Room struct {
	id      RoomID                // 房间ID
	users   map[UserID]oceanus.ID // 用户列表
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

	var list []oceanus.ID
	for _, id := range room.users {
		list = append(list, id)
	}

	room.process.SendByIDs(list, msg)
}
