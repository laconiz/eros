package message

type Encoder interface {
	Encode(msg interface{}) (*Message, error)
	Decode(stream []byte) (*Message, error)
}

type Message struct {
	Meta    Meta
	Msg     interface{}
	Raw     []byte
	Stream  []byte
	Encoder Encoder
}
