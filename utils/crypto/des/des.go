package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

func New(key []byte) (*DES, error) {

	blk, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &DES{blk: blk, iv: key}, nil
}

type DES struct {
	iv  []byte
	blk cipher.Block
}

func (d *DES) Encrypt(src []byte) []byte {

	enc := cipher.NewCBCEncrypter(d.blk, d.iv)

	s := enc.BlockSize()
	if p := s - len(src)%s; p != 0 {
		pd := bytes.Repeat([]byte{byte(p)}, p)
		src = append(src, pd...)
	}

	dst := make([]byte, len(src))
	enc.CryptBlocks(dst, src)
	return dst
}

func (d *DES) Decrypt(dst []byte) []byte {

	dec := cipher.NewCBCDecrypter(d.blk, d.iv)

	src := make([]byte, len(dst))
	dec.CryptBlocks(src, dst)

	if len(src) > 0 {
		p := int(src[len(src)-1])
		if p != 0 && p < len(src) {
			src = src[:len(src)-p]
		}
	}

	return src
}

func Encrypt(src, key []byte) ([]byte, error) {

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	enc := cipher.NewCBCEncrypter(block, key)

	s := block.BlockSize()
	if p := s - len(src)%s; p != 0 {
		pd := bytes.Repeat([]byte{byte(p)}, p)
		src = append(src, pd...)
	}

	dst := make([]byte, len(src))
	enc.CryptBlocks(dst, src)
	return dst, nil
}

func Decrypt(dst, key []byte) ([]byte, error) {

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dec := cipher.NewCBCDecrypter(block, key)

	src := make([]byte, len(dst))
	dec.CryptBlocks(src, dst)

	if len(src) > 0 {
		p := int(src[len(src)-1])
		if p != 0 && p < len(src) {
			src = src[:len(src)-p]
		}
	}

	return src, nil
}
