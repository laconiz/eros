// consul监视服务

package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api/watch"
)

type param = map[string]interface{}

type Watcher struct {
	consul *Consul
}

func (w *Watcher) Keyprefix(prefix string, handler watch.HandlerFunc) (*Plan, error) {

	if handler == nil {
		return nil, fmt.Errorf("nil handler")
	}

	plan, err := watch.Parse(param{"type": "keyprefix", "prefix": prefix})
	if err != nil {
		return nil, err
	}

	plan.Handler = handler
	return &Plan{addr: w.consul.addr, plan: plan}, nil
}

type Plan struct {
	addr string
	plan *watch.Plan
}

func (p *Plan) Run() error {
	return p.plan.Run(p.addr)
}

func (p *Plan) Stop() {
	p.plan.Stop()
}
