package process

import (
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/database/consul"
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/oceanus/proto"
	"hash/fnv"
)

const prefix = "oceanus/meshes/"

func (process *Process) register() error {
	info := process.local.Info()
	return consulor.KV().Store(string(prefix+info.ID), info.Addr)
}

func (process *Process) deregister() error {
	info := process.local.Info()
	return consulor.KV().Delete(string(prefix + info.ID))
}

func (process *Process) watcher() *consul.Plan {

	handler := func(_ uint64, value interface{}) {

		pairs := value.(api.KVPairs)

		meshes := map[proto.MeshID]string{}
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
func (process *Process) syncConnections(meshes map[proto.MeshID]string) {

	process.mutex.Lock()
	defer process.mutex.Unlock()

	local := process.local.Info()

	hash := fnv.New32()
	hash.Write([]byte(local.Addr))
	lp := hash.Sum32()

	for id, addr := range meshes {

		if id == local.ID {
			continue
		}

		if _, ok := process.connectors[id]; ok {
			continue
		}

		hash.Reset()
		hash.Write([]byte(addr))
		rp := hash.Sum32()

		if lp > rp && (lp-rp)%2 == 0 || lp < rp && (rp-lp)%2 != 0 {
			continue
		}

		connector := process.NewConnector(addr)
		process.connectors[id] = connector
		go connector.Run()
		process.logger.Data(addr).Info("connect to mesh")
	}

	for id, mesh := range process.remotes {
		if _, ok := meshes[id]; !ok {
			mesh.Destroy()
			delete(process.remotes, id)
			process.logger.Data(mesh.Info()).Info("mesh destroyed")
		}
	}

	for id, connector := range process.connectors {
		if _, ok := meshes[id]; !ok {
			connector.Stop()
			delete(process.connectors, id)
			process.logger.Data(connector.Addr()).Info("connector stopped")
		}
	}
}
