package itemMgr

import (
	"github.com/laconiz/eros/oceanus/example/proto"
	"testing"
)

const defaultUser UserID = 1000000

func TestChange(t *testing.T) {

	rds.Key().Delete(key(defaultUser))

	latest, success, err := Change(defaultUser, proto.ItemCoin, 100, proto.ItemChangeByAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if !success || latest != 100 {
		t.Fatal(latest, success)
	}

	latest, success, err = Change(defaultUser, proto.ItemCoin, -50, proto.ItemChangeByAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if !success || latest != 50 {
		t.Fatal(latest, success)
	}

	latest, success, err = Change(defaultUser, proto.ItemCoin, -100, proto.ItemChangeByAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if success || latest != 50 {
		t.Fatal(latest, success)
	}
}

func TestNum(t *testing.T) {

	latest, err := Num(defaultUser, proto.ItemCoin)
	if err != nil {
		t.Fatal(err)
	}
	if latest != 50 {
		t.Fatal(latest)
	}

	latest, err = Num(defaultUser, proto.ItemTicket)
	if err != nil {
		t.Fatal(err)
	}
	if latest != 0 {
		t.Fatal(latest)
	}
}

func TestItems(t *testing.T) {

	items, err := Items(defaultUser)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 || items[proto.ItemCoin] != 50 {
		t.Fatalf("%#v", items)
	}
}
