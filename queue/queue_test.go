package queue

import (
	"testing"
	"time"
)

func TestQueue(t *testing.T) {

	queue := NewQueue(10)

	for i := 0; i < 10; i++ {
		if queue.Add(i) != nil {
			t.FailNow()
		}
	}

	if queue.Add(0) != ErrOverflow {
		t.FailNow()
	}

	queue.Close()
	if queue.Add(0) != ErrClosed {
		t.FailNow()
	}

	queue = NewQueue(0)

	count := 0

	go func() {
		for {
			queue.Add(0)
			count++
		}
	}()

	go func() {
		for {
			queue.Pick()
		}
	}()

	time.Sleep(time.Second)
	t.Logf("1 second count: %d", count)
}
