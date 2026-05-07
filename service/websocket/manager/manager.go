package manager

import (
	"InnerG/pkg/logger"
	"sync"
)

type ConnectionManager struct {
	connections map[string]*UserConnection // 这里正常应该使用 userId + deviceId 做为 key
	mu          sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*UserConnection),
		mu:          sync.RWMutex{},
	}
}

func (m *ConnectionManager) IsConnected(ID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.connections[ID]; ok {
		return true
	}
	return false
}

func (m *ConnectionManager) GetConnection(ID string) *UserConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connections[ID]
}

func (m *ConnectionManager) AddConnection(connection *UserConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[connection.UserID] = connection
}

func (m *ConnectionManager) RemoveConnection(ID string) {
	m.mu.RLock()
	if _, ok := m.connections[ID]; !ok {
		logger.Log.Error("RemoveConnection: connection not exist")
		m.mu.Unlock()
	}
	m.mu.RUnlock()
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.connections, ID)
}
