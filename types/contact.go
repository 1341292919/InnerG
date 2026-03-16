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
