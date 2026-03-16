package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type ChatSession struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`    // MongoDB 自动生成的 ID
	SessionID string        `bson:"sessionId" json:"sessionId"` // 会话ID，唯一索引
	UserID    string        `bson:"userId" json:"userId"`       // 用户ID
	Model     string        `bson:"model" json:"model"`         // 使用的模型 (如: gpt-4)
	Title     string        `bson:"title" json:"title"`         // 会话标题
	Status    string        `bson:"status" json:"status"`       // 状态 (active/archived/deleted)
	Messages  []Message     `bson:"messages" json:"messages"`   // 消息列表
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"` // 创建时间
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"` // 更新时间
}

// Message 对应 ai_chat_sessions 中的 messages 数组元素
type Message struct {
	Role      string    `bson:"role" json:"role"`           // 角色 (system/user/assistant)
	Message   string    `bson:"message" json:"message"`     // 消息内容
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"` // 消息创建时间
}

func (m *Message) DefaultCreateAt() {
	if !m.CreatedAt.IsZero() {
		return
	}
	m.CreatedAt = time.Now()
}
