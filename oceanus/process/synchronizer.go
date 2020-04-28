package process

import (
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/remote"
	"github.com/laconiz/eros/utils/json"
)

const prefix = "oceanus/meshes/"

// 同步网格
func (proc *Process) synchronize(_ uint64, value interface{}) {

	// 数据列表
	pairs := value.(api.KVPairs)

	// 构建网格列表
	meshes := map[proto.MeshID]*proto.Mesh{}
	for _, pair := range pairs {

		info := &proto.Mesh{}
		// 非法格式
		if err := json.Unmarshal(pair.Value, info); err != nil {
			proc.logger.Err(err).Warn("invalid mesh info")
			continue
		}

		meshes[info.ID] = info
	}

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	local := proc.local.Info()

	// 同步网格列表
	for _, info := range meshes {

		// 本地网格
		if info.ID == local.ID {
			continue
		}

		// 已存在网格
		if _, ok := proc.remotes[info.ID]; ok {
			continue
		}

		// 新建网格
		proc.logger.Data(info).Info("new mesh")
		mesh := remote.New(info, proc)
		proc.remotes[info.ID] = mesh
	}

	// 清理网格列表
	for _, mesh := range proc.remotes {

		info := mesh.Info()
		id := info.ID

		// 存在网格记录
		if _, ok := meshes[id]; ok {
			continue
		}

		// 不存在网格记录 但仍存在连接
		if _, connected := mesh.State(); connected {
			continue
		}

		// 销毁网格
		mesh.Destroy()
		delete(proc.remotes, id)
		proc.logger.Data(info).Info("mesh destroyed")
	}
}
