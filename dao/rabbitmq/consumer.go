package rabbitmq

import (
	"fmt"

	"github.com/wagslane/go-rabbitmq"
)

type Consumer struct {
	consumer *rabbitmq.Consumer
	exchange string
	Queue    string
}

// Handler 定义了处理消息的业务逻辑函数类型
type Handler func(d rabbitmq.Delivery) rabbitmq.Action

func NewConsumer(exchange string, queueName string, routingKey string) (*Consumer, error) {
	// 注意：NewConsumer 的第二个参数现在是 queue name (string)
	consumer, err := rabbitmq.NewConsumer(
		conn,
		queueName, // 直接传 queue 名字，不再需要 WithConsumerOptionsQueueName
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsExchangeName(exchange),
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey), // 绑定多个路由键
	)

	if err != nil {
		return nil, fmt.Errorf("Create Consumer Failed: %w", err)
	}

	return &Consumer{consumer: consumer, exchange: exchange, Queue: queueName}, nil
}

// Run 启动监听
func (c *Consumer) Run(handler Handler) error {
	err := c.consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		return handler(d)
	})
	if err != nil {
		return fmt.Errorf("Run Consumer Failed: %w", err)
	}
	return nil
}

// Close 提供手动关闭消费者的接口
func (c *Consumer) Close() {
	if c.consumer != nil {
		c.consumer.Close()
	}
}
