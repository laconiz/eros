package cipher

type Cipher interface {
	Encipher([]byte) ([]byte, error)
	Decipher([]byte) ([]byte, error)
}

type AesCipher struct {
}
