package des

import (
	"bytes"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestDES(t *testing.T) {

	des, _ := New([]byte("12345678"))

	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {

		wg.Add(1)

		go func() {

			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for j := 0; j < 10000; j++ {

				src := make([]byte, r.Intn(20)+10)
				for k := 0; k < len(src); k++ {
					src[k] = byte(r.Int())
				}

				if !bytes.Equal(src, des.Decrypt(des.Encrypt(src))) {
					t.FailNow()
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
