package cipher

type Cipher interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

type Maker interface {
	New() Cipher
}
