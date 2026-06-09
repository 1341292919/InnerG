package websocket

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/oss"
	"context"
	"crypto/rand"
	"encoding/base64"
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

// SyncMessagesHTTP 通过HTTP接口进行游标同步（合并私聊+群聊）
func (ws *WebSocketSrv) SyncMessagesHTTP(ctx context.Context, userID, lastID int64, limit int) (*SyncResponse, error) {
	return ws.SyncMessages(ctx, userID, lastID, limit)
}

func (ws *WebSocketSrv) CreateTicket(ctx context.Context, userID int64) (string, error) {
	ticket, err := newWebSocketTicket()
	if err != nil {
		return "", err
	}
	websocketDao := dao.NewWebsocketDao()
	if err := websocketDao.Cache.SetTicket(ctx, websocketTicketKey(ticket), strconv.FormatInt(userID, 10), constants.WebsocketTicketExpire); err != nil {
		return "", err
	}
	return ticket, nil
}

func (ws *WebSocketSrv) ConsumeTicket(ctx context.Context, ticket string) (int64, error) {
	websocketDao := dao.NewWebsocketDao()
	userIDStr, err := websocketDao.Cache.ConsumeTicket(ctx, websocketTicketKey(ticket))
	if err != nil {
		return 0, err
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse websocket ticket user id: %w", err)
	}
	return userID, nil
}

func (ws *WebSocketSrv) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	if err := oss.CheckFileSize(file, constants.WebsocketImageMaxSize); err != nil {
		return "", fmt.Errorf("check file size failed: %w", err)
	}
	return ws.uploadFile(ctx, file, validateImageFile, constants.WebsocketImageFileNamePrefix, constants.WebsocketImageOssOrigin)
}

func (ws *WebSocketSrv) UploadVideo(ctx context.Context, file *multipart.FileHeader) (string, error) {
	if err := oss.CheckFileSize(file, constants.WebsocketVideoMaxSize); err != nil {
		return "", fmt.Errorf("check file size failed: %w", err)
	}
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

func websocketTicketKey(ticket string) string {
	return constants.WebsocketTicketKeyPrefix + ticket
}

func newWebSocketTicket() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "ws_" + base64.RawURLEncoding.EncodeToString(b), nil
}
