package types

type FriendTargetReq struct {
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type FriendResp struct {
	UserID    int64 `json:"user_id"`
	FriendID  int64 `json:"friend_id"`
	Status    int8  `json:"status"`
	CreatedAt int64 `json:"created_at"`
}

type FriendListResp struct {
	Friends []*FriendResp `json:"friends"`
	Total   int           `json:"total"`
}
