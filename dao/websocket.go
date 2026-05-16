package dao

import (
	"InnerG/dao/cache"
	"InnerG/dao/db"
	_interface "InnerG/dao/interface"
)

type WebsocketDao struct {
	Db    _interface.WebSocketDB
	Cache _interface.WebSocketCache
}

func NewWebsocketDao() *WebsocketDao {
	return &WebsocketDao{
		Db:    db.NewWebSocketDBClient(),
		Cache: cache.NewWebsocketRedisClient(),
	}
}
