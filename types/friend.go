package types

type FriendTargetReq struct {
	FriendID int64  `json:"friend_id" form:"friend_id" binding:"required"`
	Message  string `json:"message" form:"message"`
}

type FriendPageReq struct {
	PageSize int `form:"page_size"`
	PageNum  int `form:"page_num"`
}

type HandleFriendRequestReq struct {
	FriendID   int64  `json:"friend_id" form:"friend_id" binding:"required"`
	ActionType string `json:"action_type" form:"action_type" binding:"required"`
}

type FriendResp struct {
	UserID    int64  `json:"user_id"`
	FriendID  int64  `json:"friend_id"`
	Status    int8   `json:"status"`
	CreatedAt int64  `json:"created_at"`
	Avatar    string `json:"avatar"`
	Username  string `json:"username"`
}

type FriendListResp struct {
	Friends    []*FriendResp `json:"friends"`
	TotalCount int64         `json:"total_count"`
	PageIndex  int           `json:"page_index"`
	PageSize   int           `json:"page_size"`
}

type FriendRequestResp struct {
	ID             int64  `json:"id"`
	FromUser       int64  `json:"from_user"`
	ToUser         int64  `json:"to_user"`
	Status         int8   `json:"status"`
	Message        string `json:"message"`
	CreatedAt      int64  `json:"created_at"`
	FromUserAvatar string `json:"from_user_avatar"`
	FromUserName   string `json:"from_user_name"`
}

type FriendRequestListResp struct {
	Requests   []*FriendRequestResp `json:"requests"`
	TotalCount int64                `json:"total_count"`
	PageIndex  int                  `json:"page_index"`
	PageSize   int                  `json:"page_size"`
}
