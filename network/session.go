package network

type SessionID uint64

type Session interface {
	ID() SessionID
	Addr() string
	Send(msg interface{})
	SendStream(stream []byte)
	Close()
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{})
}

type defaultSession struct {
}

func (ses *defaultSession) ID() SessionID {
	return 0
}

func (ses *defaultSession) Addr() string {
	return ""
}

func (ses *defaultSession) Send(_ interface{}) {

}

func (ses *defaultSession) SendStream(_ []byte) {

}

func (ses *defaultSession) Close() {

}

func (ses *defaultSession) Set(key interface{}, value interface{}) {

}

func (ses *defaultSession) Get(key interface{}) (value interface{}) {
	return nil
}

var DefaultSession = &defaultSession{}
