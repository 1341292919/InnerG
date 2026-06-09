package websocket

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/logger"
	"InnerG/service/websocket/message"
	"context"
	"sort"
	"sync"
)

// SyncRequest 客户端游标同步请求
type SyncRequest struct {
	Action string `json:"action"` // "sync"
	LastID int64  `json:"last_id"`
	Limit  int    `json:"limit,omitempty"`
}

// SyncResponse 同步响应
type SyncResponse struct {
	Type     string            `json:"type"`     // "sync_messages"
	Messages []message.Message `json:"messages"` // 增量消息列表
	HasMore  bool              `json:"has_more"` // 是否还有更多
	NextID   int64             `json:"next_id"`  // 下次同步的游标
}

// SyncMessages 游标同步（合并私聊+群聊）
func (ws *WebSocketSrv) SyncMessages(ctx context.Context, userID, lastID int64, limit int) (*SyncResponse, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	var (
		privateMessages []*model.Message
		groupMessages   []*model.GroupMessage
		wg              sync.WaitGroup
		err1, err2      error
	)

	// 并发查询私聊和群聊
	wg.Add(2)

	go func() {
		defer wg.Done()
		websocketDao := dao.NewWebsocketDao()
		privateMessages, err1 = websocketDao.Db.GetMessagesAfterID(ctx, userID, lastID, limit)
	}()

	go func() {
		defer wg.Done()
		groupDao := dao.NewGroupDao()
		groupMessages, err2 = groupDao.Db.GetGroupMessagesAfterID(ctx, userID, lastID, limit)
	}()

	wg.Wait()

	if err1 != nil {
		logger.Log.Errorf("SyncMessages: get private messages error: %v", err1)
		return nil, err1
	}
	if err2 != nil {
		logger.Log.Errorf("SyncMessages: get group messages error: %v", err2)
		return nil, err2
	}

	// 合并排序（按雪花ID）
	merged := mergeAndSortMessages(privateMessages, groupMessages, limit)

	hasMore := len(merged) > limit
	if hasMore {
		merged = merged[:limit]
	}

	nextID := lastID
	if len(merged) > 0 {
		nextID = merged[len(merged)-1].ID
	}

	// 转换为业务层消息格式
	messages := make([]message.Message, len(merged))
	for i, m := range merged {
		messages[i] = message.Message{
			ID:        m.MsgID,
			UserID:    m.FromUser,
			TargetID:  m.ToUser,
			Content:   m.Content,
			Type:      m.Type,
			Status:    m.Status,
			ChatType:  m.ChatType,
			CreatedAt: m.CreatedAt,
		}
	}

	return &SyncResponse{
		Type:     "sync_messages",
		Messages: messages,
		HasMore:  hasMore,
		NextID:   nextID,
	}, nil
}

// mergedMessage 统一的消息结构（用于排序）
type mergedMessage struct {
	ID        int64
	MsgID     string
	FromUser  int64
	ToUser    int64
	Content   string
	Type      int8
	Status    int8
	ChatType  int8
	CreatedAt int64
}

// mergeAndSortMessages 合并私聊和群聊消息，按雪花ID排序
func mergeAndSortMessages(privateMessages []*model.Message, groupMessages []*model.GroupMessage, limit int) []mergedMessage {
	// 转换为统一格式
	var merged []mergedMessage

	for _, m := range privateMessages {
		merged = append(merged, mergedMessage{
			ID:        m.ID,
			MsgID:     m.MsgID,
			FromUser:  m.FromUser,
			ToUser:    m.ToUser,
			Content:   m.Content,
			Type:      m.Type,
			Status:    m.Status,
			ChatType:  constants.ChatTypePrivate,
			CreatedAt: m.CreatedAt,
		})
	}

	for _, m := range groupMessages {
		merged = append(merged, mergedMessage{
			ID:        m.ID,
			MsgID:     m.MsgID,
			FromUser:  m.FromUser,
			ToUser:    m.GroupID, // 群ID放在 ToUser 字段
			Content:   m.Content,
			Type:      m.Type,
			Status:    m.Status,
			ChatType:  constants.ChatTypeGroup,
			CreatedAt: m.CreatedAt,
		})
	}

	// 按雪花ID排序
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].ID < merged[j].ID
	})

	// 限制数量
	if len(merged) > limit+1 {
		merged = merged[:limit+1]
	}

	return merged
}
