package blocker

import (
	"sync"
	"time"
)

func New(max int, dur time.Duration) *Blocker {
	return &Blocker{
		max: max,
		rd:  dur / 2,
		tm:  time.Now(),
		pre: map[string]time.Time{},
		cur: map[string]time.Time{},
	}
}

type Blocker struct {
	max int
	rd  time.Duration
	tm  time.Time
	pre map[string]time.Time
	cur map[string]time.Time
	mu  sync.Mutex
}

func (d *Blocker) Block(s string) bool {

	d.mu.Lock()
	defer d.mu.Unlock()

	// rotate
	tm := time.Now()
	if len(d.cur) >= d.max || tm.Sub(d.tm) >= d.rd {
		d.pre = d.cur
		d.cur = map[string]time.Time{}
	}

	// verify
	if _, ok := d.pre[s]; ok {
		return true
	}
	if _, ok := d.cur[s]; ok {
		return true
	}

	// record
	d.cur[s] = tm
	return false
}
