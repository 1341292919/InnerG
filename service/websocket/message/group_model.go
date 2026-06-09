package message

import "github.com/goccy/go-json"

type GroupMessage struct {
	ID        string `json:"id"`
	UserID    int64  `json:"user_id"`
	GroupID   int64  `json:"group_id"`
	Content   string `json:"content"`
	Type      int8   `json:"type"`
	Status    int8   `json:"status"`
	ChatType  int8   `json:"chat_type"`
	CreatedAt int64  `json:"created_at"`
}

type GroupOfflineMessage struct {
	Message       GroupMessage `json:"message"`
	TargetUserIDs []int64      `json:"target_user_ids"`
}

func (m GroupMessage) JsonContent() ([]byte, error) {
	return json.Marshal(m)
}
