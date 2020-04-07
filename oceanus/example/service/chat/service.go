package chat

// ---------------------------------------------------------------------------------------------------------------------

type Service struct {
	rooms map[RoomID]*Room
	users map[UserID]*Room
}

func (service *Service) Init() {

}

func (service *Service) OnEnter(userID UserID, req *EnterREQ) {

	if room, ok := service.users[userID]; ok {

		if room.id == req.ID {
			return
		}

		service.OnQuit(userID, &QuitREQ{ID: room.id})
	}

	room, ok := service.rooms[req.ID]
	if !ok {
		return
	}

	room.users[userID] = 
}

func (service *Service) OnSpeak(userID UserID, req *SpeakREQ) {

	room, ok := service.users[userID]
	if !ok {
		return
	}

	room.Broadcast(&SpeakACK{User: userID, Content: req.content})
}

func (service *Service) OnQuit(userID UserID, req *QuitREQ) {

	room, ok := service.users[userID]
	if !ok {
		return
	}

	room.Broadcast()
}

func (service *Service) Destroy() {

}
