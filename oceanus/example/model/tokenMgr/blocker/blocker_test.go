package blocker

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestBlocker_Block(t *testing.T) {

	blocker := New(10000, time.Second)

	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {

		wg.Add(1)

		go func(i int) {

			s := strconv.Itoa(i)

			for j := 0; j < 100000000; j++ {
				blocker.Block(s)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}
