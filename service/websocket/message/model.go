package message

import "github.com/goccy/go-json"

// 私聊消息
type Message struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	TargetID  string `json:"target_id"`
	Content   string `json:"content"`
	Type      int64  `json:"type"`
	CreatedAt int64  `json:"created_at"`
}

func (m Message) JsonContent() ([]byte, error) {
	return json.Marshal(m)
}
