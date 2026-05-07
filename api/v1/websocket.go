package v1

import (
	service "InnerG/service/websocket"
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
				// 发送ping消息（心跳包）
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("发送心跳失败: %v", err)
					return
				}
			}
		}()
		l := service.NewWebSocketSrv()
		l.NewConnection(c.Request.Context(), conn)
	}
}
