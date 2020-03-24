package itemMgr

import (
	"github.com/laconiz/eros/oceanus/example/model"
	"testing"
)

const defaultUser UserID = 1000000

func TestChange(t *testing.T) {

	rds.Key().Delete(key(defaultUser))

	latest, success, err := Change(defaultUser, model.Coin, 100, model.ChangeByAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if !success || latest != 100 {
		t.Fatal(latest, success)
	}

	latest, success, err = Change(defaultUser, model.Coin, -50, model.ChangeByAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if !success || latest != 50 {
		t.Fatal(latest, success)
	}

	latest, success, err = Change(defaultUser, model.Coin, -100, model.ChangeByAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if success || latest != 50 {
		t.Fatal(latest, success)
	}
}

func TestNum(t *testing.T) {

	latest, err := Num(defaultUser, model.Coin)
	if err != nil {
		t.Fatal(err)
	}
	if latest != 50 {
		t.Fatal(latest)
	}

	latest, err = Num(defaultUser, model.Ticket)
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

	if len(items) != 1 || items[model.Coin] != 50 {
		t.Fatalf("%#v", items)
	}
}
