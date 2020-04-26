package process

import (
	"encoding/hex"
	"github.com/laconiz/eros/oceanus/proto"
	uuid "github.com/satori/go.uuid"
)

func (proc *Process) NewKeyNode(typo proto.NodeType, key proto.NodeKey, invoker interface{}) error {

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	id := proto.NodeID(hex.EncodeToString(uuid.NewV1().Bytes()))
	info := &proto.Node{ID: id, Type: typo, Key: key, Mesh: proc.local.Info().ID}

	_, err := proc.local.Insert(info, invoker)
	if err != nil {
		return err
	}

	proc.broadcast(&proto.NodeJoin{Nodes: []*proto.Node{info}})
	return nil
}

func (proc *Process) NewTypeNode(typo proto.NodeType, invoker interface{}) error {
	return proc.NewKeyNode(typo, proto.NodeKey(hex.EncodeToString(uuid.NewV1().Bytes())), invoker)
}
