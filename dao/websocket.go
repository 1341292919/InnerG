package dao

import (
	"InnerG/dao/db"
	_interface "InnerG/dao/interface"
	"context"
)

type WebsocketDao struct {
	Db _interface.WebSocketDB
}

func NewWebsocketDao(ctx context.Context) *WebsocketDao {
	return &WebsocketDao{
		Db: db.NewWebSocketDBClient(),
	}
}
