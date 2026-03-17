package types

// gin框架对require参数缺失的拦截，需要添加添加 required字段
type NewChatSessionReq struct {
	InitialMessage string `form:"initialMessage" binding:"required"`
	SessionTitle   string `form:"sessionTitle" binding:"required"`
}

type NewChatSessionResp struct {
	SessionId string
}

type StreamChatReq struct {
	SessionId   string `form:"sessionId" binding:"required"`
	UserMessage string `form:"userMessage" binding:"required"`
}
type StreamChatResp struct {
}
type GetUserSessionListReq struct {
	PageSize int `form:"pageSize" binding:"required"`
	PageNum  int `form:"pageNum" binding:"required"`
}
type GetUserSessionListResp struct {
	Total       int
	SessionList []*Session
}

type GetUserSessionDetailReq struct {
	SessionId string `form:"sessionId" binding:"required"`
}
type GetUserSessionDetailResp struct {
	SessionDetail *SessionDetail
}

type Session struct {
	UserId        string
	SessionId     string
	Model         string
	Title         string
	UpdatedAt     int64
	CreatedAt     int64
	LastMessage   string
	LastSpeakRole string
	Status        string
}

type SessionDetail struct {
	UserId    string
	SessionId string
	Model     string
	Title     string
	Messages  []Message
	Status    string
	UpdatedAt int64
	CreatedAt int64
}

type Message struct {
	Role      string
	Content   string
	CreatedAt int64
}
