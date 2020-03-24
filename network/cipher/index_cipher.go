package cipher

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

type indexCipher struct {
	sender   uint32     // 发送序列号
	receiver uint32     // 读取序列号
	rand     *rand.Rand // 随机数生成器
}

const (
	randSize = 4 // 随机数长度
	flagSize = 4 // 序列号长度
)

func cipherBytes(stream []byte) []byte {
	return []byte{
		stream[0]>>6<<6 | stream[1]>>4<<6>>2 | stream[2]>>2<<6>>4 | stream[3]<<6>>6,
		stream[1]>>6<<6 | stream[2]>>4<<6>>2 | stream[3]>>2<<6>>4 | stream[0]<<6>>6,
		stream[2]>>6<<6 | stream[3]>>4<<6>>2 | stream[0]>>2<<6>>4 | stream[1]<<6>>6,
		stream[3]>>6<<6 | stream[0]>>4<<6>>2 | stream[1]>>2<<6>>4 | stream[2]<<6>>6,
	}
}

func (c *indexCipher) Encode(raw []byte) ([]byte, error) {

	c.sender++

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, c.rand.Uint32()); err != nil {
		return nil, fmt.Errorf("write seed error: %w", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, c.sender); err != nil {
		return nil, fmt.Errorf("write flag error: %w", err)
	}
	if _, err := buf.Write(raw); err != nil {
		return nil, fmt.Errorf("write raw error: %w", err)
	}

	stream := buf.Bytes()
	cipher := cipherBytes(stream)
	for index := randSize; index < len(stream); index++ {
		stream[index] ^= cipher[index%randSize]
	}

	return stream, nil
}

func (c *indexCipher) Decode(stream []byte) ([]byte, error) {

	if len(stream) <= randSize+flagSize {
		return nil, fmt.Errorf("invalid stream size %d", len(stream))
	}

	cipher := cipherBytes(stream)
	for index := randSize; index < len(stream); index++ {
		stream[index] ^= cipher[index%randSize]
	}

	c.receiver++
	flag := binary.LittleEndian.Uint32(stream[randSize : randSize+flagSize])
	if flag != c.receiver {
		return nil, fmt.Errorf("invalid flag: %d != %d", flag, c.receiver)
	}

	return stream[randSize+flagSize:], nil
}

// ---------------------------------------------------------------------------------------------------------------------

func NewIndexMaker() Maker {
	return &indexCipherMaker{}
}

type indexCipherMaker struct {
}

func (m *indexCipherMaker) New() Cipher {
	return &indexCipher{rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}
