// session管理器

package session

import "sync"

func NewManager() *Manager {
	return &Manager{sessions: map[ID]Session{}}
}

type Manager struct {
	sessions map[ID]Session
	mutex    sync.RWMutex
}

// 获取一个session
func (mgr *Manager) Load(id ID) Session {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return mgr.sessions[id]
}

// 当前的session数量
func (mgr *Manager) Count() int64 {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return int64(len(mgr.sessions))
}

// 插入一个session
func (mgr *Manager) Insert(ses Session) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.sessions[ses.ID()] = ses
}

// 删除一个session
func (mgr *Manager) Remove(ses Session) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	delete(mgr.sessions, ses.ID())
}

// 遍历所有的session
func (mgr *Manager) Range(handler func(Session) bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	for _, ses := range mgr.sessions {
		if !handler(ses) {
			return
		}
	}
}
