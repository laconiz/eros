package redis

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestSingleton_Exec(t *testing.T) {

	assert(r.Key().Delete(key) == nil)

	var count int32
	var wg sync.WaitGroup

	f := func() {
		atomic.AddInt32(&count, 1)
	}

	for i := 0; i < 5; i++ {

		wg.Add(1)

		go func() {
			if _, err := r.Singleton(key).Exec(f); err != nil {
				t.Fatal(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	assert(count == 1)

	assert(r.Key().Delete(key) == nil)
}
