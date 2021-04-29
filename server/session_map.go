package server

import (
	"sync"
)

// SessionMap session管理
type SessionMap struct {
	sync.RWMutex
	sync.Map
}

// 存储玩家session，已经存在时，覆盖旧数据，并返回旧数据
func (m *SessionMap) StoreOrUpdate(userId int64, sess *Session) (old *Session) {
	m.Lock()
	defer m.Unlock()
	actual, loaded := m.LoadOrStore(userId, sess)
	if loaded {
		old = actual.(*Session)
		m.Store(userId, sess)
	}
	return
}

// 删除session
func (m *SessionMap) Delete(userId int64, sess *Session) {
	m.Lock()
	defer m.Unlock()
	tmp, ok := m.Load(userId)
	if !ok {
		return
	}
	old := tmp.(*Session)
	if old == sess {
		m.Map.Delete(userId)
	}
	return
}
