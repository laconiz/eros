package elastic

import (
	"github.com/laconiz/eros/utils/json"
	"testing"
)

func TestElastic_Insert(t *testing.T) {

	logs := NewPointerIncreaseLogs(1)
	if err := client.Insert(logs[0]); err != nil {
		t.Fatal(err)
	}

	err := client.Insert(nil)
	if err == nil {
		t.FailNow()
	}
	t.Log(err)

	slogs := NewStructIncreaseLogs(1)
	err = client.Insert(slogs[0])
	if err != nil {
		t.Fatal(err)
	}
}

func TestElastic_InsertRaw(t *testing.T) {

	logs := NewPointerIncreaseLogs(1)
	raw, err := json.Marshal(logs[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := client.InsertRaw("item_increase_log", raw); err != nil {
		t.Fatal(err)
	}

	err = client.InsertRaw("item_increase_log", []byte{})
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
}

func TestElastic_Inserts(t *testing.T) {

	logs := NewPointerIncreaseLogs(5)
	err := client.Inserts(logs)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Inserts(append(logs, nil))
	if err == nil {
		t.FailNow()
	}
	t.Log(err)

	slogs := NewStructIncreaseLogs(5)
	err = client.Inserts(slogs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestElastic_InsertRaws(t *testing.T) {

	logs := NewPointerIncreaseLogs(5)
	var raws [][]byte
	for _, log := range logs {
		raw, err := json.Marshal(log)
		if err != nil {
			t.Fatal(err)
		}
		raws = append(raws, raw)
	}
	err := client.InsertRaws(indexName, raws)
	if err != nil {
		t.Fatal(err)
	}

	raws = append(raws, []byte{})
	err = client.InsertRaws(indexName, raws)
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
}
