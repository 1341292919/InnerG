package v1

import (
	"InnerG/pack"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/errno"
	service "InnerG/service/websocket"
	"InnerG/types"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WebSocketConnectionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		upgrader := websocket.Upgrader{
			HandshakeTimeout: time.Second * 10,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("升级连接失败: %v", err)
			return
		}
		defer conn.Close()

		log.Println("客户端连接成功")

		// 消息体大小限制 64KB，超限自动断开
		conn.SetReadLimit(64 * 1024)

		// 设置心跳超时时间：30秒内没收到客户端消息就关闭连接
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		// 心跳重置函数，每次收到客户端消息就刷新超时时间
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			return nil
		})

		// 启动心跳发送协程，每10秒发一次心跳包
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
					log.Printf("发送心跳失败: %v", err)
					return
				}
			}
		}()
		l := service.NewWebSocketSrv()
		l.NewConnection(c.Request.Context(), conn)
	}
}

func GetWebSocketMessages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetMessagesReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		if req.PageSize <= 0 {
			req.PageSize = 20
		}
		if req.PageNum <= 0 {
			req.PageNum = 1
		}
		uid := ctl.GetUserInfo(ctx.Request.Context()).Id
		srv := service.NewWebSocketSrv()
		msgs, total, err := srv.GetMessagesByTimeRange(ctx.Request.Context(), uid, req.TargetID, req.After, req.Before, req.PageSize, req.PageNum)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.GetMessagesResp{
			Messages:   pack.BuildMessageList(msgs),
			TotalCount: total,
			PageIndex:  req.PageNum,
			PageSize:   req.PageSize,
		})
	}
}

func CreateWebSocketTicket() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctl.GetUserInfo(ctx.Request.Context()).Id
		ticket, err := service.NewWebSocketSrv().CreateTicket(ctx.Request.Context(), uid)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.CreateWebSocketTicketResp{
			Ticket:    ticket,
			ExpiresIn: int64(constants.WebsocketTicketExpire.Seconds()),
		})
	}
}

func UploadWebsocketImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, fileHeader, err := ctx.Request.FormFile("file")
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		url, err := service.NewWebSocketSrv().UploadImage(ctx.Request.Context(), fileHeader)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.WebsocketUploadResp{URL: url})
	}
}

func UploadWebsocketVideo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, fileHeader, err := ctx.Request.FormFile("file")
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		url, err := service.NewWebSocketSrv().UploadVideo(ctx.Request.Context(), fileHeader)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.WebsocketUploadResp{URL: url})
	}
}
