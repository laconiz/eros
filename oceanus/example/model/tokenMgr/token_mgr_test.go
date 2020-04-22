package tokenMgr

import (
	"github.com/laconiz/eros/oceanus/example/model"
	"testing"
)

func TestNewVerify(t *testing.T) {

	for i := 0; i < 100; i++ {

		id := model.UserID(10000000)

		token, err := New(id)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(token)

		nid, err := Verify(token)
		if err != nil {
			t.Fatal(err)
		}

		if id != nid {
			t.FailNow()
		}
	}
}
