package network

type HandlerFunc func(*Event)

type Event struct {
	Meta    *Meta
	Message interface{}
	Raw     []byte
	Stream  []byte
	Session Session
}

// ---------------------------------------------------------------------------------------------------------------------

type Connected struct {
}

type Disconnected struct {
}

type ConnectFailed struct {
}

var (
	MetaConnected     *Meta
	MetaDisconnected  *Meta
	MetaConnectFailed *Meta
)

func init() {

	RegisterMeta(Connected{}, JsonCodec)
	RegisterMeta(Disconnected{}, JsonCodec)
	RegisterMeta(ConnectFailed{}, JsonCodec)

	MetaConnected = MetaByMsg(Connected{})
	MetaDisconnected = MetaByMsg(Disconnected{})
	MetaConnectFailed = MetaByMsg(ConnectFailed{})
}
