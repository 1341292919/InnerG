package manager

import (
	"time"

	"github.com/gorilla/websocket"
)

type UserConnection struct {
	UserID     string
	Conn       *websocket.Conn
	LastActive time.Time
	DeviceID   string // 多端登录
}
