package cipher

type emptyCipher struct {
}

func (c *emptyCipher) Encode(raw []byte) ([]byte, error) {
	return raw, nil
}

func (c *emptyCipher) Decode(raw []byte) ([]byte, error) {
	return raw, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func NewEmptyMaker() Maker {
	return &emptyEncoderMaker{cipher: &emptyCipher{}}
}

type emptyEncoderMaker struct {
	cipher Cipher
}

func (m *emptyEncoderMaker) New() Cipher {
	return m.cipher
}
