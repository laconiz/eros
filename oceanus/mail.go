package oceanus

import (
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
)

type MailID string

type Mail struct {
	ID        MailID   `json:"id"`
	Senders   []NodeID `json:"senders"`
	Receivers []NodeID `json:"receivers"`
	UserID    uint64   `json:"userID"`
	MsgID     string   `json:"msgID"`
	Body      []byte   `json:"body"`
}

func (m *Mail) copy() *Mail {
	return &Mail{
		ID: MailID(hex.EncodeToString(uuid.NewV1().Bytes())),
	}
}

func (m *Mail) Receiver(receivers []NodeID) *Mail {

}
