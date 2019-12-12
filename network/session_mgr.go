package network

import (
	"sync"
	"sync/atomic"
)

type SessionMgr struct {
	sessions map[SessionID]Session
	flag     uint64
	mutex    sync.RWMutex
}

func (mgr *SessionMgr) NewID() SessionID {
	return SessionID(atomic.AddUint64(&mgr.flag, 1))
}

func (mgr *SessionMgr) Get(id SessionID) Session {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return mgr.sessions[id]
}

func (mgr *SessionMgr) Count() int64 {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return int64(len(mgr.sessions))
}

func (mgr *SessionMgr) Add(ses Session) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.sessions[ses.ID()] = ses
}

func (mgr *SessionMgr) Del(ses Session) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	delete(mgr.sessions, ses.ID())
}

func (mgr *SessionMgr) Range(inv func(Session) bool) {

	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	for _, ses := range mgr.sessions {
		if !inv(ses) {
			return
		}
	}
}

func NewSessionMgr() *SessionMgr {
	return &SessionMgr{sessions: map[SessionID]Session{}}
}
