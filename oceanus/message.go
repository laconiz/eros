package oceanus

import (
	"errors"
	"fmt"
	"github.com/laconiz/eros/message"
	uuid "github.com/satori/go.uuid"
)

type Event struct {
	UUID      string
	Type      string
	Name      string
	MessageID message.ID
	Body      []byte
	Chain     []*ChannelInfo
	channel   *Channel
	peer      *Peer
}

func (e *Event) Send(typo string, name string, msg interface{}) error {

	if e.peer == nil {
		return errors.New("event not from any peer")
	}

	meta, ok := message.MetaByMessage(msg)
	if !ok {
		return fmt.Errorf("can not get meta by message: %#v", msg)
	}

	raw, err := meta.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode message error: %w", err)
	}

	event := &Event{
		UUID:      uuid.NewV1().String(),
		Type:      typo,
		Name:      name,
		MessageID: meta.ID(),
		Body:      raw,
		Chain:     append(e.Chain, &e.channel.ChannelInfo),
	}

	return e.peer.Send(event)
}
