package network

import (
	"bytes"
	"fmt"
	"github.com/laconiz/eros/json"
	"reflect"
	"unsafe"
)

type Encoder interface {
	Encode(msg interface{}) (*Event, error)
	Decode(stream []byte) (*Event, error)
	String(*Event) string
}

// ---------------------------------------------------------------------------------------------------------------------

type nameEncoder struct {
	sep byte
}

func (e *nameEncoder) Encode(msg interface{}) (*Event, error) {

	meta := MetaByMsg(msg)
	if meta == nil {
		return nil, fmt.Errorf("non-proto message: %+v", msg)
	}

	raw, err := meta.Codec().Encode(msg)
	if err != nil {
		return nil, fmt.Errorf("encode message %+v error: %w", msg, err)
	}

	buf := bytes.NewBufferString(meta.Name())
	buf.WriteByte(e.sep)
	buf.Write(raw)

	return &Event{Meta: meta, Msg: msg, Raw: raw, Stream: buf.Bytes()}, nil
}

var messageIDSize = int(unsafe.Sizeof(MessageID(0)))

func (e *nameEncoder) Decode(stream []byte) (*Event, error) {

	if len(stream) < messageIDSize {
		return nil, fmt.Errorf("invalid stream length: %d", len(stream))
	}

	buf := bytes.NewBuffer(stream)

	name, err := buf.ReadString(e.sep)
	if err != nil {
		return nil, fmt.Errorf("read message name error: %w", err)
	}
	if len(name) > 0 {
		name = name[:len(name)-1]
	}

	meta := MetaByName(name)
	if meta == nil {
		return nil, fmt.Errorf("invalid mesage name: %s", name)
	}

	raw := buf.Bytes()

	msg := reflect.New(meta.Type()).Interface()
	if err := meta.Codec().Decode(raw, msg); err != nil {
		return nil, fmt.Errorf("decode message %s error: %w", string(stream), err)
	}

	return &Event{Meta: meta, Msg: msg, Raw: raw, Stream: stream}, nil
}

func (e *nameEncoder) String(event *Event) string {

	if event == nil {
		return "nil event"
	}

	meta := event.Meta
	if event.Meta == nil {
		meta = MetaByMsg(event.Msg)
	}
	if meta == nil {
		return "invalid event"
	}

	if event.Meta.Codec() == JsonCodec && event.Stream != nil {
		return string(event.Stream)
	}

	raw, err := json.Marshal(event.Msg)
	if err != nil {
		return fmt.Sprintf("marshal message[%#v] to json error: %v", event.Msg, err)
	}

	return fmt.Sprintf("%s%s%s", meta.Name(), string(e.sep), string(raw))
}

var NameEncoder Encoder = &nameEncoder{sep: '-'}
