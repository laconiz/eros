package process

import (
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/database/consul"
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/oceanus/proto"
)

func (process *Process) register() error {
	info := process.local.Info()
	return consulor.KV().Store(string(prefix+info.ID), info)
}

func (process *Process) deregister() error {
	info := process.local.Info()
	return consulor.KV().Delete(string(prefix + info.ID))
}

func (process *Process) watcher() *consul.Plan {

	handler := func(_ uint64, value interface{}) {

		pairs := value.(api.KVPairs)

		meshes := map[string]*proto.Mesh{}
		err := consul.ParsePairs(prefix, pairs, &meshes, false)
		if err != nil {
			process.logger.Err(err).Error("parse meshes error")
			return
		}

		process.syncConnections(meshes)
	}

	plan, _ := consulor.Watcher().Keyprefix(prefix, handler)
	return plan
}

// 同步连接信息
func (process *Process) syncConnections(meshes map[string]*proto.Mesh) {

	process.mutex.Lock()
	defer process.mutex.Unlock()

	local := process.local.Info()

	for id, mesh := range meshes {

		if proto.MeshID(id) == local.ID {
			continue
		}

		if _, ok := process.connectors[mesh.ID]; ok {
			continue
		}

		if (local.Power > mesh.Power && (local.Power-mesh.Power)%2 == 0) ||
			(local.Power < mesh.Power && (mesh.Power-local.Power)%2 != 0) {
			continue
		}

		connector := process.NewConnector(mesh.Addr)
		process.connectors[mesh.ID] = connector
		go connector.Run()
		process.logger.Data(mesh).Info("sync mesh")
	}

	for id, mesh := range process.remotes {
		if _, ok := meshes[string(id)]; !ok {
			mesh.Destroy()
			delete(process.remotes, id)
			process.logger.Data(mesh.Info()).Info("mesh destroy")
		}
	}

	for id, connector := range process.connectors {
		if _, ok := meshes[string(id)]; !ok {
			connector.Stop()
			delete(process.connectors, id)
			process.logger.Data(id).Info("connector stopped")
		}
	}
}

const prefix = "oceanus/"
