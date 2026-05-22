package manager

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type UserConnection struct {
	UserID     int64
	Conn       *websocket.Conn
	LastActive time.Time
	DeviceID   string // 多端登录
	writeMu    sync.Mutex
}

func (uc *UserConnection) WriteData(data []byte) error {
	uc.writeMu.Lock()
	defer uc.writeMu.Unlock()
	return uc.Conn.WriteMessage(websocket.TextMessage, data)
}

func (uc *UserConnection) WriteJSONData(v any) error {
	uc.writeMu.Lock()
	defer uc.writeMu.Unlock()
	return uc.Conn.WriteJSON(v)
}
