package types

type CreateGroupReq struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description string  `json:"description" binding:"max=500"`
	MemberIDs   []int64 `json:"member_ids" binding:"required,min=1"`
}

// CreateGroupResp 创建群组响应
type CreateGroupResp struct {
	GroupID     int64  `json:"group_id"`
	Name        string `json:"name"`
	OwnerID     int64  `json:"owner_id"`
	MemberCount int    `json:"member_count"`
	CreatedAt   int64  `json:"created_at"`
}

type GetGroupsResp struct {
	Groups []*GroupInfo `json:"groups"`
	Total  int          `json:"total"`
}

type GroupInfo struct {
	GroupID     int64             `json:"group_id"`
	Name        string            `json:"name"`
	Avatar      string            `json:"avatar"`
	OwnerID     int64             `json:"owner_id"`
	Description string            `json:"description"`
	MemberCount int               `json:"member_count"`
	MaxMembers  int               `json:"max_members"`
	CreatedAt   int64             `json:"created_at"`
	LastMessage *GroupLastMessage `json:"last_message,omitempty"`
}

type GroupLastMessage struct {
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type GetGroupDetailReq struct {
	GroupID int64 `form:"group_id" binding:"required"`
}

type UpdateGroupReq struct {
	GroupID     int64   `json:"group_id" binding:"required"`
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Avatar      *string `json:"avatar" binding:"omitempty,max=512"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

type DeleteGroupReq struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type AddGroupMembersReq struct {
	GroupID int64   `json:"group_id" binding:"required"`
	UserIDs []int64 `json:"user_ids" binding:"required,min=1"`
}

type RemoveGroupMemberReq struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type QuitGroupReq struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type GetGroupMembersReq struct {
	GroupID  int64 `form:"group_id" binding:"required"`
	Page     int   `form:"page"`
	PageSize int   `form:"page_size"`
}

type GetGroupMembersResp struct {
	Members   []*GroupMemberInfo `json:"members"`
	Total     int64              `json:"total"`
	PageIndex int                `json:"page_index"`
	PageSize  int                `json:"page_size"`
}

type GroupMemberInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Role     int8   `json:"role"`
	JoinedAt int64  `json:"joined_at"`
}

type GetGroupMessagesReq struct {
	GroupID  int64 `form:"group_id" binding:"required"`
	Before   int64 `form:"before"`
	PageSize int   `form:"page_size"`
	PageNum  int   `form:"page_num"`
}

type GetGroupMessagesResp struct {
	Messages   []*GroupMessageResp `json:"messages"`
	TotalCount int64               `json:"total_count"`
	PageIndex  int                 `json:"page_index"`
	PageSize   int                 `json:"page_size"`
}

type GroupMessageResp struct {
	MsgID      string `json:"msg_id"`
	FromUser   int64  `json:"from_user"`
	FromName   string `json:"from_name"`
	FromAvatar string `json:"from_avatar"`
	Content    string `json:"content"`
	Type       int8   `json:"type"`
	CreatedAt  int64  `json:"created_at"`
}
