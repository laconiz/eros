package oceanus

import (
	"fmt"
	"math/rand"
	"time"
)

type Mesh interface {
	Info() *Node
	Push(*Message) error
	Connected() bool
}

type Route interface {
	Info() *Channel
	Node() Mesh
	Push(*Message) error
}

type Router struct {

	// 类型
	typo string

	// 脏数据标记
	dirty bool

	// 通道数据
	routes map[string]Route

	// 平衡节点
	balances []Route

	// 随机数
	rand *rand.Rand
}

// 将通道数据置为过期
func (r *Router) expired() {
	r.dirty = true
}

// 添加一个通道
func (r *Router) add(route Route) {
	r.expired()
	r.routes[route.Info().ID] = route
}

// 删除一个通道
func (r *Router) remove(route Route) {
	r.expired()
	delete(r.routes, route.Info().ID)
}

func (r *Router) rotate() {

	for _, route := range r.routes {

		if route.Node().Connected() {

		}

	}

	r.dirty = false
}

func (r *Router) Balance(message *Message) error {

	if r.dirty {
		r.rotate()
	}

	if len(r.balances) == 0 {
		return fmt.Errorf("can not find channel on type %v", r.typo)
	}

	index := r.rand.Intn(len(r.balances))
	return r.balances[index].Push(message)
}

func (r *Router) Send(message *Message) error {

	return nil
}

func NewRouter(typo string) *Router {
	return &Router{
		typo:   typo,
		dirty:  false,
		routes: map[string]Route{},
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
