package types

type GetMessagesReq struct {
	TargetID int64 `form:"target_id" binding:"required"`
	After    int64 `form:"after"`
	Before   int64 `form:"before"`
	PageSize int   `form:"page_size"`
	PageNum  int   `form:"page_num"`
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
	Messages   []*MessageResp `json:"messages"`
	TotalCount int64          `json:"total_count"`
	PageIndex  int            `json:"page_index"`
	PageSize   int            `json:"page_size"`
}

type GetUnreadResp struct {
	Messages []*MessageResp `json:"messages"`
	Total    int            `json:"total"`
}
