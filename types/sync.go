package types

import "InnerG/service/websocket/message"

// SyncMessagesReq 游标同步请求
type SyncMessagesReq struct {
	LastID int64 `form:"last_id" binding:"required"` // 最后一条消息的雪花ID
	Limit  int   `form:"limit"`                      // 每次同步数量，默认100，最大200
}

// SyncMessagesResp 游标同步响应
type SyncMessagesResp struct {
	Messages []message.Message `json:"messages"` // 增量消息列表
	HasMore  bool              `json:"has_more"` // 是否还有更多
	NextID   int64             `json:"next_id"`  // 下次同步的游标
}
