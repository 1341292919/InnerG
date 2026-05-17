package websocket

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"context"
)

func (ws *WebSocketSrv) GetMessagesAfterTimestamp(ctx context.Context, userID, targetID, after int64) ([]*model.Message, error) {
	websocketDao := dao.NewWebsocketDao()
	return websocketDao.Db.GetMessagesAfterTimestamp(ctx, userID, targetID, after)
}

func (ws *WebSocketSrv) GetUnreadMessages(ctx context.Context, userID int64) ([]*model.Message, error) {
	websocketDao := dao.NewWebsocketDao()
	return websocketDao.Db.GetOfflineMessages(ctx, userID, constants.OnceOffMessagePushNum)
}
