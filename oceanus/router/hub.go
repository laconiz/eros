package router

import "github.com/laconiz/eros/oceanus/proto"

type Hub struct {
	typo proto.NodeType
}

func (h *Hub) Type() proto.NodeType {
	return h.typo
}
