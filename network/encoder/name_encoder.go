// 消息流编码规则

package encoder

import (
	"bytes"
	"errors"
	"github.com/laconiz/eros/network/message"
)

// ---------------------------------------------------------------------------------------------------------------------

type nameEncoder struct {
}

func (e *nameEncoder) Marshal(msg interface{}) (*message.Message, error) {

	meta, ok := message.MetaByMsg(msg)
	if !ok {
		return nil, errors.New("meta cannot be found")
	}

	raw, err := meta.Encode(msg)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString(meta.Name())
	buf.WriteByte(NameEncoderSep())
	buf.Write(raw)

	return &message.Message{Meta: meta, Msg: msg, Raw: raw, Stream: buf.Bytes()}, nil
}

func (e *nameEncoder) Unmarshal(stream []byte) (*message.Message, error) {

	bp := bytes.SplitN(stream, []byte{NameEncoderSep()}, 2)
	if len(bp) != 2 {
		return nil, errors.New("invalid stream")
	}

	meta, ok := message.MetaByName(string(bp[0]))
	if !ok {
		return nil, errors.New("meta cannot be found")
	}

	msg, err := meta.Decode(bp[1])
	if err != nil {
		return nil, err
	}

	return &message.Message{Meta: meta, Msg: msg, Raw: bp[1], Stream: stream}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func NameEncoderSep() byte {
	return '-'
}

// ---------------------------------------------------------------------------------------------------------------------

func NewNameMaker() Maker {
	return &nameEncoderMaker{encoder: &nameEncoder{}}
}

type nameEncoderMaker struct {
	encoder Encoder
}

func (m *nameEncoderMaker) New() Encoder {
	return m.encoder
}
