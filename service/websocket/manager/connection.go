package manager

import (
	"time"

	"github.com/gorilla/websocket"
)

type UserConnection struct {
	UserID     int64
	Conn       *websocket.Conn
	LastActive time.Time
	DeviceID   string // 多端登录
}
