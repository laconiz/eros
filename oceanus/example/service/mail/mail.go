package mail

import (
	"github.com/laconiz/eros/oceanus/example/model"
	"github.com/laconiz/eros/oceanus/example/proto"
)

func New() *Service {
	return &Service{users: map[model.UserID]*model.Mail{}}
}

type Service struct {
	users  map[model.UserID]*model.Mail
	global []*model.Mail
}

func (service *Service) Init() {

}

func (service *Service) OnMails(req *proto.MailListREQ) *proto.MailListACK {

}

func (service *Service) OnNewMail() {

}

func (service *Service) OnReadMail() {

}

func (service *Service) OnReceiveMail(id proto.UserID) {

}

func (service *Service) Destroy() {

}
