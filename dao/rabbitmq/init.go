package rabbitmq

import (
	"InnerG/pkg/utils"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

var conn *rabbitmq.Conn

func Init() {
	var err error

	conn, err = rabbitmq.NewConn(
		utils.WithRabbitMqUrl(),
		rabbitmq.WithConnectionOptionsLogging, // 启用日志
		rabbitmq.WithConnectionOptionsConfig(rabbitmq.Config{
			Heartbeat: 60 * time.Second, // 设置心跳间隔
		}),
		rabbitmq.WithConnectionOptionsReconnectInterval(5*time.Second), // 重连间隔
	)
	if err != nil {
		panic(err)
	}
}
