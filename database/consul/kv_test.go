package consul

import (
	"github.com/laconiz/eros/utils/json"
	"testing"
)

type Value struct {
	Int    int
	String string
}

func testLoad(t *testing.T, key string, value interface{}) {
	t.Logf("test load key %v", key)
	if err := consul.KV().Load(key, value); err != nil {
		t.Fatalf("load key %v error: %v", key, err)
	}
	t.Logf("load {%v:%v} success", key, json.String(value))
}

func testStore(t *testing.T, key string, value interface{}) {
	t.Logf("test store {%v:%v}", key, json.String(value))
	if err := consul.KV().Store(key, value); err != nil {
		t.Fatalf("store {%v:%v} error: %v", key, json.String(value), err)
	}
	t.Logf("store {%v:%v} success", key, json.String(value))
}

func testDelete(t *testing.T, key string) {
	t.Logf("test delete key %v", key)
	if err := consul.KV().Delete(key); err != nil {
		t.Fatalf("delete key %v error: %v", key, err)
	}
	if err := consul.KV().Load(key, &Value{}); err != ErrNotFound {
		t.Fatalf("delete key %v error: %v", key, nil)
	}
	t.Logf("delete key %v success", key)
}

func testLoadsStrict(t *testing.T, prefix string, value interface{}) {
	t.Logf("test load prefix %v", prefix)
	err := consul.KV().Loads(prefix, value, true)
	if err == nil {
		t.Fatalf("load prefix %v error: %v", prefix, err)
	}
	t.Logf("load {%v:%v} success: %v", prefix, json.String(value), err)
}

func testLoadsNoStrict(t *testing.T, prefix string, value interface{}) {
	t.Logf("test load prefix %v", prefix)
	if err := consul.KV().Loads(prefix, value, false); err != nil {
		t.Fatalf("load prefix %v error: %v", prefix, err)
	}
	t.Logf("load {%v:%v} success", prefix, json.String(value))
}

func TestKV(t *testing.T) {

	prefix := "test/"

	suffixA := "A"
	suffixC := "C"
	suffixD := "D"

	keyA := prefix + suffixA
	keyC := prefix + suffixC
	keyD := prefix + suffixD

	valueA := Value{Int: 1, String: "africa"}
	valueB := Value{}
	valueC := Value{Int: 3, String: "china"}
	valueD := "different"

	testStore(t, keyA, valueA)
	testLoad(t, keyA, &valueB)
	testStore(t, keyC, &valueC)
	testStore(t, keyD, valueD)

	structValues := map[string]Value{}
	testLoadsStrict(t, prefix, &structValues)
	structValues = map[string]Value{}
	testLoadsNoStrict(t, prefix, &structValues)
	if len(structValues) != 2 {
		t.Fatalf("load prefix %v value count error: need 2, got %d", prefix, len(structValues))
	}
	if structValues[suffixA] != valueA {
		t.Fatalf("load prefix %v value error: %v != %v",
			prefix, json.String(structValues[suffixA]), json.String(valueA))
	}
	if structValues[suffixC] != valueC {
		t.Fatalf("load prefix %v value error: %v != %v",
			prefix, json.String(structValues[suffixC]), json.String(valueC))
	}

	pointerValues := map[string]*Value{}
	testLoadsStrict(t, prefix, &pointerValues)
	pointerValues = map[string]*Value{}
	testLoadsNoStrict(t, prefix, &pointerValues)
	if len(pointerValues) != 2 {
		t.Fatalf("load prefix %v value count error: need 2, got %d", prefix, len(pointerValues))
	}
	if *pointerValues[suffixA] != valueA {
		t.Fatalf("load prefix %v value error: %v != %v",
			prefix, json.String(pointerValues[suffixA]), json.String(valueA))
	}
	if *pointerValues[suffixC] != valueC {
		t.Fatalf("load prefix %v value error: %v != %v",
			prefix, json.String(pointerValues[suffixC]), json.String(valueC))
	}

	testDelete(t, keyA)
	testDelete(t, keyC)
	testDelete(t, keyD)
}
