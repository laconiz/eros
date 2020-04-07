package reader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// ---------------------------------------------------------------------------------------------------------------------

type sizeReader struct {
}

func (p *sizeReader) Write(conn net.Conn, stream []byte) error {

	size := int32(len(stream))
	if size == 0 {
		return errors.New("send empty stream")
	}

	if err := binary.Write(conn, binary.LittleEndian, size); err != nil {
		return fmt.Errorf("write header error: %w", err)
	}

	if _, err := conn.Write(stream); err != nil {
		return fmt.Errorf("write stream error: %w", err)
	}

	return nil
}

func (p *sizeReader) Read(conn net.Conn) ([]byte, error) {

	var size int32
	if err := binary.Read(conn, binary.LittleEndian, &size); err != nil {
		return nil, fmt.Errorf("read header error: %w", err)
	}

	stream := make([]byte, size)
	if _, err := conn.Read(stream); err != nil {
		return nil, fmt.Errorf("read stream error: %w", err)
	}

	return stream, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func NewSizeMaker() Maker {
	return &sizeReaderMaker{packer: &sizeReader{}}
}

type sizeReaderMaker struct {
	packer Reader
}

func (m *sizeReaderMaker) New() Reader {
	return m.packer
}
