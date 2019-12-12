package network

type State string

const (
	Stopped State = "stopped"
	Running State = "running"
	Closing State = "closing"
)

type Service interface {
	Start()
	Stop()
	State() State
}

type Acceptor interface {
	Service
	Count() int64
}

type Connector interface {
	Service
}
