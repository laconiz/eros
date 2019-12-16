package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/laconiz/eros/network"
	"net"
	"reflect"
	"unsafe"
)

type Encoder interface {
	Encode(net.Conn) (*network.Event, error)
	Decode(*network.Event) ([]byte, error)
}

type StdEncoder struct {
}

func (enc *StdEncoder) Encode(conn net.Conn) (*network.Event, error) {

	// 读取流长度
	var size int32
	if err := binary.Read(conn, binary.LittleEndian, &size); err != nil {
		return nil, fmt.Errorf("read size error: %w", err)
	}

	// 获取流
	stream := make([]byte, size)
	n, err := conn.Read(stream)
	if err != nil {
		return nil, fmt.Errorf("read raw error: %w", err)
	}
	if int32(n) != size {
		return nil, fmt.Errorf("read raw error: need %d bytes, got %d", size, n)
	}

	reader := bytes.NewBuffer(stream)

	// 读取消息ID
	var id network.MessageID
	if err := binary.Read(reader, binary.LittleEndian, &id); err != nil {
		return nil, fmt.Errorf("read message id error: %w", err)
	}

	// 获取消息元数据
	meta := network.MetaByID(id)
	if meta == nil {
		return nil, fmt.Errorf("invalid meta id: %d", id)
	}

	// 反序列化消息
	raw := reader.Bytes()
	msg := reflect.New(meta.Type()).Interface()
	if err := meta.Codec().Decode(raw, msg); err != nil {
		return nil, fmt.Errorf("decode message[%s] by raw[%s] error: %w", meta.Type(), string(raw), err)
	}

	return &network.Event{
		Meta:   meta,
		Msg:    msg,
		Raw:    raw,
		Stream: stream,
	}, nil
}

func (enc *StdEncoder) Decode(event *network.Event) ([]byte, error) {

	if event.Stream == nil {

		// 获取消息元数据
		meta := network.MetaByMsg(event.Msg)
		if meta == nil {
			return nil, fmt.Errorf("invalid message: %#v", event.Msg)
		}

		// 序列化消息
		raw, err := meta.Codec().Encode(event.Msg)
		if err != nil {
			return nil, fmt.Errorf("encode message[%#v] error: %w", event.Msg, err)
		}
		event.Raw = raw

		// 写入消息ID
		buf := &bytes.Buffer{}
		if err := binary.Write(buf, binary.LittleEndian, meta.ID()); err != nil {
			return nil, fmt.Errorf("write message id[%d] error: %w", meta.ID(), err)
		}

		// 写入消息体
		n, err := buf.Write(raw)
		if err != nil {
			return nil, fmt.Errorf("write raw[%s] error: %w", string(raw), err)
		}
		if n != len(raw) {
			return nil, fmt.Errorf("write raw[%s] error: has %d bytes, %d wrote", string(raw), len(raw), n)
		}

		event.Stream = buf.Bytes()

	} else {

		msgSize := int(unsafe.Sizeof(network.MessageID(0)))
		if len(event.Stream) <= msgSize {
			return nil, fmt.Errorf("invalid stream size: %d", len(event.Stream))
		}

		event.Raw = event.Stream[msgSize:]
	}

	buf := &bytes.Buffer{}

	// 写入流长度
	size := int32(len(event.Stream))
	if err := binary.Write(buf, binary.LittleEndian, size); err != nil {
		return nil, fmt.Errorf("write size[%d] error: %w", size, err)
	}

	// 写入消息流
	n, err := buf.Write(event.Stream)
	if err != nil {
		return nil, fmt.Errorf("write stream[%s] error: %w", string(event.Stream), err)
	}
	if int32(n) != size {
		return nil, fmt.Errorf("write stream[%s] error: has %d bytes, %d wrote", string(event.Stream), size, n)
	}

	return buf.Bytes(), nil
}

func (enc *StdEncoder) String() string {

}

type EncoderMaker interface {
	New() Encoder
}

type StdEncoderMaker struct {
}

func (maker *StdEncoderMaker) New() Encoder {
	return &StdEncoder{}
}
