package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

func Watch() *watcher {
	return &watcher{}
}

type watcher struct {
	plan *watch.Plan
}

func (w *watcher) Run() error {
	return w.plan.Run(addr)
}

func (w *watcher) Stop() {
	w.plan.Stop()
}

func NewKeyPrefixWatcher(prefix string, handler func(api.KVPairs)) (*watcher, error) {

	plan, err := watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": prefix,
	})
	if err != nil {
		return nil, err
	}

	plan.Handler = func(_ uint64, value interface{}) {
		if pairs, ok := value.(api.KVPairs); ok {
			handler(pairs)
		}
	}
	return &watcher{plan: plan}, nil
}
