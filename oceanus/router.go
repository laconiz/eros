package oceanus

type Mesh interface {
	Info() *Mesh
	Push(*Message) error
	Connected() bool
}

type Node interface {
	Info() *Node
	Mesh() Mesh
	Push(*Message) error
}

type Router struct {
	nodes     map[NodeID]Node
	balancers map[NodeType]*Balancer
}

// 添加一个节点
func (r *Router) Add(node Node) *Balancer {

	return nil
}

func (r *Router) Remove(id NodeID) {

}

func (r *Router) Expired(typo NodeType) {

}

func NewRouter() *Router {
	return &Router{
		nodes: map[NodeID]Node{},
	}
}
