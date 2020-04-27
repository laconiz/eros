package network

type State string

const (
	Stopped State = "stopped"
	Running State = "running"
	Closing State = "closing"
)

type Service interface {
	Run()
	Stop()
}

type Acceptor interface {
	Service
	Count() int64
	State() State
	Broadcast(interface{})
	BroadcastRaw([]byte)
}

type Connector interface {
	Service
	Connected() bool
	Send(interface{}) error
	SendRaw([]byte) error
}
