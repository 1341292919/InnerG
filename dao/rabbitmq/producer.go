package rabbitmq

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/wagslane/go-rabbitmq"
)

type Producer struct {
	publisher *rabbitmq.Publisher
	exchange  string // 标注调用方
}

func NewProducer(exchange string) (*Producer, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchange),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeDurable,       // 持久化交换器
		rabbitmq.WithPublisherOptionsExchangeKind("topic"), // 使用 topic 类型
	)
	if err != nil {
		return nil, fmt.Errorf("创建生产者失败: %w", err)
	}

	return &Producer{
		publisher: publisher,
		exchange:  exchange,
	}, nil
}

func (p *Producer) SendMessage(topic string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.publisher.Publish(
		body,
		[]string{topic},
		rabbitmq.WithPublishOptionsExchange(p.exchange),
	)
}
