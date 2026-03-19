package service

import (
	"InnerG/config"
	"InnerG/dao"
	MongoModel "InnerG/dao/mongo/model"
	"InnerG/pack"
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
	"log"
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

// NewChatSession 新增一个会话
func (svc *ContactSrv) NewChatSession(ctx context.Context, req *types.NewChatSessionReq) (string, error) {
	u := ctl.GetUserInfo(ctx)
	dao := dao.NewContactDao(ctx)
	sessionId := strconv.FormatInt(rand.Int63(), 10)
	err := dao.Mongo.NewChatSession(ctx, &MongoModel.ChatSession{
		UserID:    u.Id,
		SessionID: sessionId,
		Model:     config.Api.Model,
		Status:    constants.CommonHealthStatus,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Messages:  []MongoModel.Message{},
	})
	if err != nil {
		return "", err
	}
	return sessionId, nil
}

// StreamChat 聊天
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
	message = append(message, utils.Message{Role: constants.CommonSystemRole, Content: constants.CommonOpener})
	message = append(message, convertMessageList(chatHistory.Messages)...)
	message = append(message, utils.Message{Role: constants.CommonUserRole, Content: req.UserMessage})

	shouldParseTitle := len(message) == 2
	resp, err := utils.SendMessageToAPI(message)
	if err != nil {
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	fullMessage := make([]string, 0)
	// 采集title
	done := false
	title := ""
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
			if e != nil {
				if e == io.EOF {
					break loop
				}
				return errno.InternalServiceError.WithMessage("streamTrans error pause " + e.Error())
			}
			streamResp := ParseStreamLine(line)
			if streamResp == nil {
				continue
			}
			if shouldParseTitle &&
				streamResp.Data.Content == constants.TileMarker {
				done = true
				fmt.Fprint(ctx.Writer, ConvertSSE(&types.StreamResp{
					Code:    "10000",
					Message: "completed",
				}))
				flusher.Flush()
			}
			if shouldParseTitle &&
				strings.Contains(streamResp.Data.Content, constants.TileMarker) {
				/*
					两种情况
					titleMarker前有内容 需流式写入
					titleMarker后有内容 需处理出来作为标题部分
					titleMarker前后都有内容
				*/
				markerIdx := strings.Index(streamResp.Data.Content, constants.TileMarker)
				beforeMarker := streamResp.Data.Content[:markerIdx]
				afterMarker := streamResp.Data.Content[markerIdx+len(constants.TileMarker):]

				// 标记前有内容，需要转发给前端
				if beforeMarker != "" {
					streamResp.Data.Content = beforeMarker
					fmt.Fprint(ctx.Writer, ConvertSSE(streamResp))
					fmt.Fprint(ctx.Writer, ConvertSSE(&types.StreamResp{
						Code:    "10000",
						Message: "completed",
					}))
					flusher.Flush()
					fullMessage = append(fullMessage, beforeMarker)
				}
				// 标记后有内容，需要处理出来作为标题部分
				if afterMarker != "" {
					title += afterMarker
				}
				// 标记已出现，后续内容都作为标题
				done = true
				continue
			}
			if done {
				title += streamResp.Data.Content
				continue
			}
			fullMessage = append(fullMessage, streamResp.Data.Content)
			// 直接写入前端（已是SSE格式）
			fmt.Fprint(ctx.Writer, ConvertSSE(streamResp))
			flusher.Flush()
		}
	}

	if shouldParseTitle {
		pack.WithTitle(ctx, title)
		go func() {
			err = dao.Mongo.UpdateSessionTitle(context.Background(), chatHistory.SessionID, title)
			if err != nil {
				log.Println("update session title error: " + err.Error())
			}
		}()
	}

	return dao.Mongo.InsertMessageToSession(ctx.Request.Context(), chatHistory.SessionID,
		[]MongoModel.Message{
			{
				Role:      constants.CommonUserRole,
				Message:   req.UserMessage,
				CreatedAt: time.Now(),
			},
			{
				Role:      constants.CommonBotRole,
				Message:   strings.Join(fullMessage, ""),
				CreatedAt: time.Now(),
			},
		})
}

func (svc *ContactSrv) GetUserSessionHistory(ctx context.Context) ([]*types.Session, int, error) {
	u := ctl.GetUserInfo(ctx)
	dao := dao.NewContactDao(ctx)
	sessionList, total, err := dao.Mongo.GetSessionByUserId(ctx, u.Id)
	if err != nil {
		return nil, -1, err
	}
	return pack.BuildSessionList(sessionList), total, nil
}
func (svc *ContactSrv) GetUserSessionDetail(ctx context.Context, req *types.GetUserSessionDetailReq) (*types.SessionDetail, error) {
	u := ctl.GetUserInfo(ctx)
	dao := dao.NewContactDao(ctx)
	exist, data, err := dao.Mongo.IsQuerySessionExist(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("session not exist")
	}
	if u.Id != data.UserID {
		return nil, fmt.Errorf("session not avaliable")
	}
	return pack.BuildSessionDetail(data), nil
}

// ParseStreamLine 解析SSE行数据，返回统一的StreamResp结构
func ParseStreamLine(line string) *types.StreamResp {
	// 1. 去掉 "data: " 前缀
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "data: ") {
		return nil
	}

	jsonStr := strings.TrimPrefix(line, "data: ")

	// 2. 如果是心跳或特殊事件，直接返回空
	if jsonStr == "[DONE]" {
		return &types.StreamResp{
			Code:    "10000",
			Message: "completed",
		}
	}

	// 3. 解析API返回的JSON
	var apiResp struct {
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &apiResp); err != nil {
		return nil
	}

	// 4. 提取content
	content := ""
	if len(apiResp.Choices) > 0 {
		content = apiResp.Choices[0].Delta.Content
	}

	// 5. 返回统一的StreamResp
	return &types.StreamResp{
		Code:    "10000",
		Message: "success",
		Data: types.StreamRespContent{
			Content: content,
		},
	}
}

func ConvertSSE(streamResp *types.StreamResp) string {
	jsonData, _ := json.Marshal(streamResp)
	return fmt.Sprintf("data:%s\n\n", jsonData)
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
