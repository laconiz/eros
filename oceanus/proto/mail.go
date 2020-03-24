package proto

import (
	"encoding/hex"
	"github.com/laconiz/eros/network/message"
	uuid "github.com/satori/go.uuid"
)

type MailID string

type Mail struct {
	ID        MailID
	Type      NodeType
	Senders   []Node
	Receivers []*Node
	Body      []byte
}

func (mail *Mail) New() *Mail {
	return &Mail{
		ID:        newMailID(),
		Type:      mail.Type,
		Senders:   mail.Senders,
		Receivers: mail.Receivers,
		Body:      mail.Body,
	}
}

func (mail *Mail) Copy() *Mail {
	return &Mail{
		ID:        mail.ID,
		Type:      mail.Type,
		Senders:   mail.Senders,
		Receivers: mail.Receivers,
		Body:      mail.Body,
	}
}

func newMailID() MailID {
	return MailID(hex.EncodeToString(uuid.NewV1().Bytes()))
}

func init() {
	message.Register(Mail{}, message.JsonCodec())
}
