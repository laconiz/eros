package oceanus

import (
	"github.com/laconiz/eros/network"
	"sync"
)

type Process struct {
	local *LMesh
	mutex sync.RWMutex
}

// ---------------------------------------------------------------------------------------------------------------------

func (proc *Process) onMail(event *network.Event) {

	mail := event.Msg.(*Mail)

	// RPC RESPONSE
	if mail.Reply != emptyRpcID {

		proc.mutex.Lock()
		defer proc.mutex.Unlock()

		return
	}

	// TODO PROXY MAIL
	if false {

		return
	}

	// 本地消息
	proc.mutex.RLock()
	defer proc.mutex.RUnlock()

	proc.local.Mail(mail)
}

// ---------------------------------------------------------------------------------------------------------------------

func (proc *Process) onMeshJoin(event *network.Event) {

	req := event.Msg.(*MeshJoin)

	proc.mutex.Lock()
	defer proc.mutex.Unlock()
}
