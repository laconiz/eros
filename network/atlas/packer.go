package atlas

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type Packer interface {
	UnPack(net.Conn) ([]byte, error)
	Pack([]byte) ([]byte, error)
}

type SizePacker struct {
}

func (p *SizePacker) Read(conn net.Conn) ([]byte, error) {

	var size int32
	if err := binary.Read(conn, binary.LittleEndian, &size); err != nil {
		return nil, fmt.Errorf("read size error: %w", err)
	}

	stream := make([]byte, size)
	length, err := conn.Read(stream)
	if err != nil {
		return nil, fmt.Errorf("read content error: %w", err)
	} else if length != int(size) {
		return nil, fmt.Errorf("read content length error: need %d, got %d", size, length)
	}
	return stream, nil
}

func (p *SizePacker) Write(conn net.Conn, stream []byte) error {

	size := int32(len(stream))
	if size == 0 {
		return errors.New("send empty stream")
	}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, size); err != nil {
		return fmt.Errorf("write header error: %w", err)
	}
	n, err := buf.Write(stream)
	if err != nil {
		return fmt.Errorf("write content error: %w", err)
	}
	if n != len(stream) {
		return fmt.Errorf("write content error: length %d, wrote %d", len(stream), n)
	}

	n, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("write to conn error: %w", err)
	}
	if n != buf.Len() {
		return fmt.Errorf("write to conn error: length %d, wrote %d", buf.Len(), n)
	}

	return nil
}
