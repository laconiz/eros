// 消息流编码规则

package message

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Encoder interface {
	Encode(msg interface{}) (Message, error)
	Decode(stream []byte) (Message, error)
}

// ---------------------------------------------------------------------------------------------------------------------
// 使用消息名编码消息流

type nameEncoder struct {
	MetaMgr
}

func (e *nameEncoder) Encode(msg interface{}) (Message, error) {

	meta, ok := e.MetaByMessage(msg)
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

	return &message{meta: meta, msg: msg, raw: raw, stream: buf.Bytes()}, nil
}

func (e *nameEncoder) Decode(stream []byte) (Message, error) {

	bp := bytes.SplitN(stream, []byte{NameEncoderSep()}, 1)
	if len(bp) != 2 {
		return nil, errors.New("invalid stream")
	}

	meta, ok := e.MetaByName(string(bp[0]))
	if !ok {
		return nil, errors.New("meta cannot be found")
	}

	msg, err := meta.Decode(bp[1])
	if err != nil {
		return nil, err
	}

	return &message{meta: meta, msg: msg, raw: bp[1], stream: stream}, nil
}

func NameEncoderSep() byte {
	return '-'
}

func DefaultNameEncoder() Encoder {
	return globalNameEncoder
}

func NewNameEncoder(mgr MetaMgr) Encoder {
	return &nameEncoder{MetaMgr: mgr}
}

// ---------------------------------------------------------------------------------------------------------------------
// 使用消息ID编码消息流

type idEncoder struct {
	MetaMgr
}

func (e *idEncoder) Encode(msg interface{}) (Message, error) {

	meta, ok := e.MetaByMessage(msg)
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

	return &message{meta: meta, msg: msg, raw: raw, stream: buf.Bytes()}, nil
}

func (e *idEncoder) Decode(stream []byte) (Message, error) {

	buf := bytes.NewBuffer(stream)

	var id ID
	if err := binary.Read(buf, binary.LittleEndian, &id); err != nil {
		return nil, err
	}

	meta, ok := e.MetaByID(id)
	if !ok {
		return nil, errors.New("meta cannot be found")
	}

	raw := buf.Bytes()

	msg, err := meta.Decode(raw)
	if err != nil {
		return nil, err
	}

	return &message{meta: meta, msg: msg, raw: raw, stream: stream}, nil
}

func DefaultIDEncoder() Encoder {
	return globalIDEncoder
}

func NewIDEncoder(mgr MetaMgr) Encoder {
	return &idEncoder{MetaMgr: mgr}
}

var (
	globalIDEncoder   Encoder
	globalNameEncoder Encoder
)

func init() {
	globalIDEncoder = NewIDEncoder(globalMetaMgr)
	globalNameEncoder = NewNameEncoder(globalMetaMgr)
}
