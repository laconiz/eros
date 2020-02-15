package cipher

import (
	"bytes"
	"log"
	"testing"
)

func TestIndexCipher(t *testing.T) {

	cipher := NewIndexMaker().New()

	origin := []byte("hello world")

	for i := 0; i < 100; i++ {

		stream, err := cipher.Encode(origin)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("index: %d, stream: %v", i, stream)

		raw, err := cipher.Decode(stream)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(origin, raw) {
			t.Fatalf("%v != %v", origin, raw)
		}

		// log.Printf("index: %d, stream: %v", i, stream)
	}
}
