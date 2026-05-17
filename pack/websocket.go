package pack

import (
	"InnerG/dao/db/model"
	"InnerG/types"
)

func BuildMessageList(msgs []*model.Message) []*types.MessageResp {
	resp := make([]*types.MessageResp, 0, len(msgs))
	for _, m := range msgs {
		resp = append(resp, &types.MessageResp{
			ID:        m.MsgID,
			FromUser:  m.FromUser,
			ToUser:    m.ToUser,
			Content:   m.Content,
			Type:      m.Type,
			Status:    m.Status,
			CreatedAt: m.CreatedAt,
		})
	}
	return resp
}
