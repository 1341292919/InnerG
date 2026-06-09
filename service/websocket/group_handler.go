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

	msg.Status = constants.MessagePushedStatus
	onlineSet := ws.broadcastToOnlineMembers(msg)

	if err := ws.sendOfflineMessages(ctx, msg, onlineSet); err != nil {
		return err
	}

	if err = ws.sender.SendMessage(constants.GroupStoreMessageTopic, msg); err != nil {
		logger.Log.Error("RouteGroupMessage:GroupStoreMessageTopic: ", err)
		return err
	}

	return nil
}

func (ws *WebSocketSrv) broadcastToOnlineMembers(msg *message.GroupMessage) map[int64]struct{} {
	onlineMembers := ws.groupManager.GetOnlineMembers(msg.GroupID)
	onlineSet := make(map[int64]struct{}, len(onlineMembers)+1)
	onlineSet[msg.UserID] = struct{}{}

	for _, memberID := range onlineMembers {
		onlineSet[memberID] = struct{}{}
		if memberID == msg.UserID {
			continue
		}
		conn := ws.manager.GetConnection(ws.manager.WithConnectionId(memberID))
		if conn == nil {
			continue
		}
		if err := conn.WriteJSONData(msg); err != nil {
			logger.Log.Errorf("broadcastToOnlineMembers:WriteJSONData: user_id=%d err=%v", memberID, err)
		}
	}
	return onlineSet
}

func (ws *WebSocketSrv) sendOfflineMessages(ctx context.Context, msg *message.GroupMessage, onlineSet map[int64]struct{}) error {
	allMembers, err := ws.getGroupMembersWithCache(ctx, msg.GroupID)
	if err != nil {
		logger.Log.Errorf("sendOfflineMessages:getGroupMembersWithCache: %v", err)
		return fmt.Errorf("failed to get group members")
	}

	offlineMembers := make([]int64, 0, len(allMembers))
	for _, memberID := range allMembers {
		if _, ok := onlineSet[memberID]; !ok {
			offlineMembers = append(offlineMembers, memberID)
		}
	}

	if len(offlineMembers) == 0 {
		return nil
	}

	msg.Status = constants.MessageUnPushedStatus
	if err = ws.sender.SendMessage(constants.GroupOfflineMessageTopic, message.GroupOfflineMessage{
		Message:       *msg,
		TargetUserIDs: offlineMembers,
	}); err != nil {
		logger.Log.Error("sendOfflineMessages:GroupOfflineMessageTopic: ", err)
		return err
	}
	msg.Status = constants.MessagePushedStatus
	return nil
}

func (ws *WebSocketSrv) getGroupMembersWithCache(ctx context.Context, groupID int64) ([]int64, error) {
	groupDao := dao.NewGroupDao()

	cached, err := groupDao.Cache.GetGroupMembers(ctx, groupID)
	if err == nil && len(cached) > 0 {
		return cached, nil
	}

	members, err := groupDao.Db.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}

	_ = groupDao.Cache.SetGroupMembers(ctx, groupID, members, constants.GroupMembersCacheTTL)
	return members, nil
}

func (ws *WebSocketSrv) BroadcastGroupSystemMessage(ctx context.Context, groupID int64, content string) error {
	groupDao := dao.NewGroupDao()
	now := time.Now().Unix()

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

	if err := groupDao.Db.InsertGroupMessage(ctx, &model.GroupMessage{
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
