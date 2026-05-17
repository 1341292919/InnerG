package types

type GetMessagesReq struct {
	TargetID int64 `form:"target_id" binding:"required"`
	After    int64 `form:"after" binding:"required"`
}

type MessageResp struct {
	ID        string `json:"id"`
	FromUser  int64  `json:"from_user"`
	ToUser    int64  `json:"to_user"`
	Content   string `json:"content"`
	Type      int8   `json:"type"`
	Status    int8   `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type GetMessagesResp struct {
	Messages []*MessageResp `json:"messages"`
}

type GetUnreadResp struct {
	Messages []*MessageResp `json:"messages"`
	Total    int            `json:"total"`
}
