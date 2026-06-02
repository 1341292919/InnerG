package websocket

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"context"
)

func (ws *WebSocketSrv) GetMessagesByTimeRange(ctx context.Context, userID, targetID, after, before int64, pageSize, pageNum int) ([]*model.Message, int64, error) {
	if err := canChat(ctx, userID, targetID); err != nil {
		return nil, 0, err
	}
	websocketDao := dao.NewWebsocketDao()
	return websocketDao.Db.GetMessagesByTimeRange(ctx, userID, targetID, after, before, pageSize, pageNum)
}

func (ws *WebSocketSrv) GetUnreadMessages(ctx context.Context, userID int64) ([]*model.Message, error) {
	websocketDao := dao.NewWebsocketDao()
	return websocketDao.Db.GetOfflineMessages(ctx, userID, constants.OnceOffMessagePushNum)
}
