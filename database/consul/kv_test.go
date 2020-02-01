package consul

import "testing"

type Value struct {
	A int
	B string
}

func TestKV(t *testing.T) {

	va := &Value{A: 1, B: "hello"}
	if err := KV().Store(prefix+"A", va); err != nil {
		t.Fatal(err)
	}

	vb := &Value{}
	if err := KV().Load(prefix+"A", vb); err != nil {
		t.Fatal(err)
	}

	if vb.A != va.A || vb.B != va.B {
		t.FailNow()
	}

	vc := &Value{A: 3, B: "world"}
	if err := KV().Store(prefix+"C", vc); err != nil {
		t.Fatal(err)
	}

	la := map[string]Value{}
	if err := KV().Loads(prefix, &la, false); err != nil {
		t.Fatal(err)
	}

	if len(la) != 2 ||
		la["A"].A != va.A || la["A"].B != va.B ||
		la["C"].A != vc.A || la["C"].B != vc.B {
		t.FailNow()
	}

	if err := KV().Delete(prefix + "A"); err != nil {
		t.Fatal(err)
	}
	if err := KV().Delete(prefix + "C"); err != nil {
		t.Fatal(err)
	}

	if err := KV().Load(prefix+"A", vb); err != nil {
		t.Fatal(err)
	}
}

const prefix = "test/"
