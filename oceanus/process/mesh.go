package process

import (
	"encoding/hex"
	"github.com/laconiz/eros/oceanus/proto"
	uuid "github.com/satori/go.uuid"
)

func (process *Process) NewKeyNode(typo proto.NodeType, key proto.NodeKey, invoker interface{}) error {

	process.mutex.Lock()
	defer process.mutex.Unlock()

	id := proto.NodeID(hex.EncodeToString(uuid.NewV1().Bytes()))
	info := &proto.Node{ID: id, Type: typo, Key: key, Mesh: process.local.Info().ID}

	_, err := process.local.Insert(info, invoker)
	if err != nil {
		return err
	}

	process.broadcast(&proto.NodeJoin{Nodes: []*proto.Node{info}})
	return nil
}

func (process *Process) NewTypeNode(typo proto.NodeType, invoker interface{}) error {
	return process.NewKeyNode(typo, proto.NodeKey(hex.EncodeToString(uuid.NewV1().Bytes())), invoker)
}
