package websocket

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/oss"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"
)

var (
	validateImageFile = oss.IsImage
	validateVideoFile = oss.IsVideo
	saveUploadFile    = oss.SaveFile
	uploadLocalFile   = oss.Upload
	nowUnix           = func() int64 { return time.Now().Unix() }
	getIMBucket       = func() string { return config.Oss.IMBucket }
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

// AckMessages ack消息，避免连接波动导致消息实际上未被接收到
func (ws *WebSocketSrv) AckMessages(ctx context.Context, userID int64, msgIDs []string) error {
	websocketDao := dao.NewWebsocketDao()
	return websocketDao.Db.AckMessages(ctx, userID, msgIDs)
}

func (ws *WebSocketSrv) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return ws.uploadFile(ctx, file, validateImageFile, constants.WebsocketImageFileNamePrefix, constants.WebsocketImageOssOrigin)
}

func (ws *WebSocketSrv) UploadVideo(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return ws.uploadFile(ctx, file, validateVideoFile, constants.WebsocketVideoFileNamePrefix, constants.WebsocketVideoOssOrigin)
}

func (ws *WebSocketSrv) uploadFile(ctx context.Context, file *multipart.FileHeader, validate func(*multipart.FileHeader) error, prefix, origin string) (string, error) {
	uid := ctl.GetUserInfo(ctx).Id
	if err := validate(file); err != nil {
		return "", fmt.Errorf("check upload file failed: %w", err)
	}
	fileName := fmt.Sprintf("%s_%d_%d", prefix, uid, nowUnix())
	if err := saveUploadFile(file, constants.StorePath, fileName); err != nil {
		return "", fmt.Errorf("save file failed: %w", err)
	}
	filePath := filepath.Join(constants.StorePath, fileName)
	return uploadLocalFile(filePath, fileName, strconv.FormatInt(uid, 10), origin, getIMBucket())
}
