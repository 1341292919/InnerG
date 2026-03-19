package pack

import (
	"InnerG/dao/mongo/model"
	"InnerG/types"
)

func BuildSessionList(session []*model.ChatSession) []*types.Session {
	buildSession := func(session *model.ChatSession) *types.Session {
		if len(session.Messages) == 0 {
			return &types.Session{
				SessionId:     session.SessionID,
				Title:         session.Title,
				UserId:        session.UserID,
				Model:         session.Model,
				Status:        session.Status,
				LastMessage:   "暂无信息噢",
				LastSpeakRole: "",
				CreatedAt:     session.CreatedAt.Unix(),
				UpdatedAt:     session.UpdatedAt.Unix(),
			}
		}
		return &types.Session{
			SessionId:     session.SessionID,
			Title:         session.Title,
			UserId:        session.UserID,
			Model:         session.Model,
			Status:        session.Status,
			LastMessage:   session.Messages[len(session.Messages)-1].Message,
			LastSpeakRole: session.Messages[len(session.Messages)-1].Role,
			CreatedAt:     session.CreatedAt.Unix(),
			UpdatedAt:     session.UpdatedAt.Unix(),
		}
	}
	res := make([]*types.Session, 0)
	for _, s := range session {
		res = append(res, buildSession(s))
	}
	return res
}

func BuildSessionDetail(session *model.ChatSession) *types.SessionDetail {
	return &types.SessionDetail{
		SessionId:  session.SessionID,
		Title:      session.Title,
		UserId:     session.UserID,
		Model:      session.Model,
		Status:     session.Status,
		Messages:   buildMessageList(session.Messages),
		MessageNum: len(session.Messages),
		CreatedAt:  session.CreatedAt.Unix(),
		UpdatedAt:  session.UpdatedAt.Unix(),
	}
}
func buildMessageList(message []model.Message) []types.Message {
	buildMessage := func(message model.Message) types.Message {
		return types.Message{
			Role:      message.Role,
			Content:   message.Message,
			CreatedAt: message.CreatedAt.Unix(),
		}
	}
	res := make([]types.Message, 0)
	for _, m := range message {
		res = append(res, buildMessage(m))
	}
	return res
}
