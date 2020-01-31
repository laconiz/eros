package steropes

import "testing"

func TestAcceptor(t *testing.T) {
	acceptor, err := NewAcceptor(AcceptorOption{})
	if err != nil {
		t.FailNow()
	}
	acceptor.Run()
}
