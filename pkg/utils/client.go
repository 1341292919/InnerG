package utils

import (
	"InnerG/config"
	"fmt"
)

func WithRabbitMqUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s/",
		config.RabbitMq.Username,
		config.RabbitMq.Password,
		config.RabbitMq.Addr,
	)
}
