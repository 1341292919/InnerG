package websocket

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/logger"
	"InnerG/service/websocket/message"
	"context"
	"fmt"
	"time"
)

func (ws *WebSocketSrv) RouteGroupMessage(ctx context.Context, msg *message.GroupMessage) error {
	groupDao := dao.NewGroupDao()

	isMember, err := groupDao.Db.IsMember(ctx, msg.GroupID, msg.UserID)
	if err != nil {
		logger.Log.Errorf("RouteGroupMessage:IsMember: %v", err)
		return fmt.Errorf("failed to check group membership")
	}
	if !isMember {
		logger.Log.Errorf("RouteGroupMessage:not group member: user_id=%d group_id=%d", msg.UserID, msg.GroupID)
		return fmt.Errorf("not a group member")
	}

	if !constants.IsValidGroupMessageType(msg.Type) {
		return fmt.Errorf("invalid message type: %d", msg.Type)
	}

	// 1. 先持久化消息
	msg.Status = constants.GroupMessageStatusNormal
	if err = ws.sender.SendMessage(constants.GroupStoreMessageTopic, msg); err != nil {
		logger.Log.Error("RouteGroupMessage:GroupStoreMessageTopic: ", err)
		return err
	}

	// 2. 尽力而为地广播给在线成员
	ws.broadcastToOnlineMembers(msg)

	// 不再发送离线消息到队列，客户端重连时通过游标同步获取
	return nil
}

func (ws *WebSocketSrv) broadcastToOnlineMembers(msg *message.GroupMessage) {
	onlineMembers := ws.groupManager.GetOnlineMembers(msg.GroupID)

	for _, memberID := range onlineMembers {
		if memberID == msg.UserID {
			continue
		}
		conn := ws.manager.GetConnection(ws.manager.WithConnectionId(memberID))
		if conn == nil {
			continue
		}
		if err := conn.WriteJSONData(msg); err != nil {
			logger.Log.Errorf("broadcastToOnlineMembers:WriteJSONData: user_id=%d err=%v", memberID, err)
			// 推送失败，客户端会通过游标同步补齐
		}
	}
}

func (ws *WebSocketSrv) BroadcastGroupSystemMessage(ctx context.Context, groupID int64, content string) error {
	now := time.Now().Unix()

	// 生成雪花ID
	snowflakeID, err := GenerateMessageID()
	if err != nil {
		logger.Log.Errorf("BroadcastGroupSystemMessage: generate snowflake ID error: %v", err)
		snowflakeID = now // 降级使用时间戳
	}

	sysMsg := &message.GroupMessage{
		ID:        fmt.Sprintf("sys_%d", now),
		UserID:    0,
		GroupID:   groupID,
		Content:   content,
		Type:      constants.GroupMessageTypeSystem,
		Status:    constants.GroupMessageStatusNormal,
		ChatType:  constants.ChatTypeGroup,
		CreatedAt: now,
	}

	onlineMembers := ws.groupManager.GetOnlineMembers(groupID)
	for _, memberID := range onlineMembers {
		conn := ws.manager.GetConnection(ws.manager.WithConnectionId(memberID))
		if conn != nil {
			_ = conn.WriteJSONData(sysMsg)
		}
	}

	groupDao := dao.NewGroupDao()
	if err := groupDao.Db.InsertGroupMessage(ctx, &model.GroupMessage{
		ID:        snowflakeID,
		MsgID:     sysMsg.ID,
		GroupID:   groupID,
		FromUser:  0,
		Content:   content,
		Type:      constants.GroupMessageTypeSystem,
		Status:    constants.GroupMessageStatusNormal,
		CreatedAt: sysMsg.CreatedAt,
	}); err != nil {
		logger.Log.Errorf("BroadcastGroupSystemMessage:InsertGroupMessage: %v", err)
	}

	return nil
}
