package encoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/laconiz/eros/holder/message"
)

type idEncoder struct {
}

func (e *idEncoder) Encode(msg interface{}) (*message.Message, error) {

	meta, ok := message.MetaByMessage(msg)
	if !ok {
		return nil, errors.New("meta cannot be found")
	}

	raw, err := meta.Encode(msg)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, meta.ID())
	buf.Write(raw)

	return &message.Message{Meta: meta, Msg: msg, Raw: raw, Stream: buf.Bytes(), Encoder: e}, nil
}

func (e *idEncoder) Decode(stream []byte) (*message.Message, error) {

	buf := bytes.NewBuffer(stream)

	var id message.ID
	if err := binary.Read(buf, binary.LittleEndian, &id); err != nil {
		return nil, err
	}

	meta, ok := message.MetaByID(id)
	if !ok {
		return nil, errors.New("meta cannot be found")
	}

	raw := buf.Bytes()

	msg, err := meta.Decode(raw)
	if err != nil {
		return nil, err
	}

	return &message.Message{Meta: meta, Msg: msg, Raw: raw, Stream: stream, Encoder: e}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func NewIDMaker() Maker {
	return &idEncoderMaker{encoder: &idEncoder{}}
}

type idEncoderMaker struct {
	encoder message.Encoder
}

func (m *idEncoderMaker) New() message.Encoder {
	return m.encoder
}
