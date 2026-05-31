package manager

import (
	"hash/fnv"
	"strconv"
	"sync"
)

const defaultShardNum = 32

type connShard struct {
	mu          sync.RWMutex
	connections map[string]*UserConnection
}

type ConnectionManager struct {
	shards []*connShard
}

func NewConnectionManager(shardNum int) *ConnectionManager {
	if shardNum <= 0 {
		shardNum = defaultShardNum
	}
	shards := make([]*connShard, shardNum)
	for i := range shards {
		shards[i] = &connShard{
			connections: make(map[string]*UserConnection),
		}
	}
	return &ConnectionManager{shards: shards}
}

func (m *ConnectionManager) getShard(id string) *connShard {
	h := fnv.New32a()
	h.Write([]byte(id))
	return m.shards[h.Sum32()%uint32(len(m.shards))]
}

func (m *ConnectionManager) IsConnected(ID string) bool {
	s := m.getShard(ID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.connections[ID]
	return ok
}

func (m *ConnectionManager) GetConnection(ID string) *UserConnection {
	s := m.getShard(ID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connections[ID]
}

func (m *ConnectionManager) AddConnection(connection *UserConnection) {
	id := m.WithConnectionId(connection.UserID)
	s := m.getShard(id)
	s.mu.Lock()
	defer s.mu.Unlock()
	if old, ok := s.connections[id]; ok {
		old.Conn.Close()
	}
	s.connections[id] = connection
}

func (m *ConnectionManager) RemoveConnection(connection *UserConnection) {
	id := m.WithConnectionId(connection.UserID)
	s := m.getShard(id)
	s.mu.Lock()
	defer s.mu.Unlock()
	if current, ok := s.connections[id]; ok && current == connection {
		delete(s.connections, id)
	}
}

func (m *ConnectionManager) WithConnectionId(userId int64) string {
	return strconv.FormatInt(userId, 10)
}
