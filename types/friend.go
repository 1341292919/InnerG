package types

type FriendTargetReq struct {
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type HandleFriendRequestReq struct {
	FriendID   int64  `json:"friend_id" form:"friend_id" binding:"required"`
	ActionType string `json:"action_type" form:"action_type" binding:"required"`
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

type FriendRequestResp struct {
	ID        int64 `json:"id"`
	FromUser  int64 `json:"from_user"`
	ToUser    int64 `json:"to_user"`
	Status    int8  `json:"status"`
	CreatedAt int64 `json:"created_at"`
}

type FriendRequestListResp struct {
	Requests []*FriendRequestResp `json:"requests"`
	Total    int                  `json:"total"`
}
