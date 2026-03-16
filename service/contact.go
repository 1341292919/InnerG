package service

import (
	"InnerG/config"
	"InnerG/dao"
	MongoModel "InnerG/dao/mongo/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/errno"
	"InnerG/pkg/utils"
	"InnerG/types"
	"bufio"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ContactSrvIns *ContactSrv
var ContactSrvOnce sync.Once

type ContactSrv struct{} // 空结构体，只包含方法

func GetContactSrv() *ContactSrv {
	ContactSrvOnce.Do(func() {
		ContactSrvIns = &ContactSrv{}
	})
	return ContactSrvIns
}

// NewContact 需要流式传输，传入初始ctx
func (svc *ContactSrv) NewContact(ctx *gin.Context) error {
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("get http flusher error")
	}
	message := make([]utils.Message, 0)
	message = append(message, utils.Message{
		Role:    constants.CommonUserRole,
		Content: constants.CommonOpener,
	})
	resp, err := utils.SendMessageToAPI(message)
	if err != nil {
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	fullMessage := make([]string, 0)

loop:
	for {
		var line string
		var e error
		select {
		case <-ctx.Request.Context().Done():
			resp.Body.Close()
			break loop
		default:
			line, e = reader.ReadString('\n')
			fullMessage = append(fullMessage, line)
			if e != nil {
				if e == io.EOF {
					break loop
				}
				return errno.InternalServiceError.WithMessage("streamTrans error pause " + e.Error())
			}
			// 直接写入前端（已是SSE格式）
			fmt.Fprint(ctx.Writer, line)
			flusher.Flush()
		}
	}
	// 这边需要解析出resp的回复 并将其存进数据库
	u := ctl.GetUserInfo(ctx.Request.Context())
	dao := dao.NewContactDao(ctx.Request.Context())
	err = dao.Mongo.NewChatSession(ctx.Request.Context(), &MongoModel.ChatSession{
		UserID:    u.Id,
		SessionID: strconv.FormatInt(rand.Int63(), 10),
		Model:     config.Api.Model,
		Title:     "ningning",
		Status:    constants.CommonHealthStatus,
		Messages: []MongoModel.Message{
			{
				Role:      constants.CommonUserRole,
				Message:   constants.CommonOpener,
				CreatedAt: time.Now(),
			},
			{
				Role:      constants.CommonBotRole,
				Message:   ParseRespContent(fullMessage),
				CreatedAt: time.Now(),
			},
		},
	})
	return err
}

// NewChatSession 新增一个会话
func (svc *ContactSrv) NewChatSession(ctx context.Context, req *types.NewChatSessionReq) (string, error) {
	u := ctl.GetUserInfo(ctx)
	dao := dao.NewContactDao(ctx)
	sessionId := strconv.FormatInt(rand.Int63(), 10)
	err := dao.Mongo.NewChatSession(ctx, &MongoModel.ChatSession{
		UserID:    u.Id,
		SessionID: sessionId,
		Model:     config.Api.Model,
		Title:     req.SessionTitle,
		Status:    constants.CommonHealthStatus,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Messages: []MongoModel.Message{
			{
				Role:      constants.CommonUserRole,
				Message:   req.InitialMessage,
				CreatedAt: time.Now(),
			},
		},
	})
	if err != nil {
		return "", err
	}
	return sessionId, nil
}

func (svc *ContactSrv) StreamChat(ctx *gin.Context, req *types.StreamChatReq) error {
	// 查询会话记录
	dao := dao.NewContactDao(ctx)
	exist, chatHistory, err := dao.Mongo.IsQuerySessionExist(ctx.Request.Context(), req.SessionId)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("session not exist")
	}
	u := ctl.GetUserInfo(ctx.Request.Context())
	if u.Id != chatHistory.UserID {
	}
	// 添加会话上文
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("get http flusher error")
	}
	message := make([]utils.Message, 0)
	message = append(message, convertMessageList(chatHistory.Messages)...)
	message = append(message, utils.Message{Role: constants.CommonUserRole, Content: req.UserMessage})

	resp, err := utils.SendMessageToAPI(message)
	if err != nil {
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	fullMessage := make([]string, 0)

loop:
	for {
		var line string
		var e error
		select {
		case <-ctx.Request.Context().Done():
			resp.Body.Close()
			break loop
		default:
			line, e = reader.ReadString('\n')
			fullMessage = append(fullMessage, line)
			if e != nil {
				if e == io.EOF {
					break loop
				}
				return errno.InternalServiceError.WithMessage("streamTrans error pause " + e.Error())
			}
			// 直接写入前端（已是SSE格式）
			fmt.Fprint(ctx.Writer, line)
			flusher.Flush()
		}
	}

	return dao.Mongo.InsertMessageToSession(ctx.Request.Context(), chatHistory.SessionID,
		[]MongoModel.Message{
			{
				Role:    constants.CommonUserRole,
				Message: req.UserMessage,
			},
			{
				Role:    constants.CommonBotRole,
				Message: ParseRespContent(fullMessage),
			},
		})
}
func ParseRespContent(fullMessage []string) string {
	var result strings.Builder

	for _, line := range fullMessage {
		// 跳过空行和 [DONE] 标记
		if line == "" || line == "data: [DONE]" {
			continue
		}

		// 移除 "data: " 前缀
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		// 解析 JSON
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			// 如果解析失败，跳过这一条
			continue
		}

		// 提取 content 内容
		if len(chunk.Choices) > 0 {
			result.WriteString(chunk.Choices[0].Delta.Content)
		}
	}

	return result.String()
}

func convertMessageList(mes []MongoModel.Message) []utils.Message {
	convertMessage := func(message MongoModel.Message) utils.Message {
		return utils.Message{
			Role:    message.Role,
			Content: message.Message,
		}
	}
	res := make([]utils.Message, 0)
	for _, m := range mes {
		res = append(res, convertMessage(m))
	}
	return res
}
