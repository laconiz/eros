package gateway

import (
	"bytes"
	"fmt"
	"github.com/laconiz/eros/iris"
	"github.com/laconiz/eros/network"
)

type encoder struct {
}

func (e *encoder) Encode(msg interface{}) (*network.Event, error) {
	return iris.Encoder.Encode(msg)
}

func (e *encoder) Decode(stream []byte) (*network.Event, error) {

	buf := bytes.NewBuffer(stream)

	// 读取消息包名
	raw, err := buf.ReadBytes('.')
	if err != nil {
		return nil, fmt.Errorf("read package name error: %w", err)
	}
	// 去除分隔符
	if len(raw) > 0 {
		raw = raw[:len(raw)-1]
	}

	// 包名
	pkg := string(raw)
	switch pkg {

	case "gateway":
		// 网关所属消息解析
		return iris.Encoder.Decode(stream)

	default:
		// 路由消息
		return &network.Event{
			Msg:    nil,
			Stream: stream,
		}, nil
	}
}

func (e *encoder) String(event *network.Event) string {
	return string(event.Stream)
}
