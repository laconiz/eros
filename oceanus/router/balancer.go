package router

import (
	"container/list"
	"github.com/laconiz/eros/oceanus/proto"
	"sort"
)

func balance(nodes map[Key]Node) *list.List {

	hash := map[proto.MeshID]*Credit{}

	for _, node := range nodes {

		mesh := node.Mesh()
		id := mesh.Info().ID

		// 查询网格结构
		credit, ok := hash[id]
		// 创建网格结构
		if !ok {

			// 查询网格状态
			state, ok := mesh.State()
			// 网格未连接
			if !ok {
				continue
			}

			// 创建网格结构
			credit = &Credit{
				mesh:  mesh,
				state: state,
			}
			hash[id] = credit
		}

		// 插入节点
		credit.nodes = append(credit.nodes, node)
	}

	// 构造排序列表
	var credits Credits
	for _, credit := range hash {
		credits = append(credits, credit)
	}
	// 排序
	sort.Sort(credits)

	// 构造链表
	list := list.New()
	for _, credit := range credits {
		for _, node := range credit.nodes {
			list.PushBack(node)
		}
	}
	return list
}

// 均衡记录
type Credit struct {
	mesh  Mesh         // 网格
	state *proto.State // 状态
	nodes []Node       // 节点列表
}

// 均衡列表
type Credits []*Credit

func (cred Credits) Len() int {
	return len(cred)
}

func (cred Credits) Less(i, j int) bool {
	ic := cred[i].state.Credit
	jc := cred[j].state.Credit
	return ic <= jc
}

func (cred Credits) Swap(i, j int) {
	cred[i], cred[j] = cred[j], cred[i]
}
