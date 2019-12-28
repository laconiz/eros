package codec

type Codec interface {
	Encode(msg interface{}) (raw []byte, err error)
	Decode(raw []byte, msg interface{}) (err error)
}
