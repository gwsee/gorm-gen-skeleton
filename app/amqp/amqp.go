package amqp

import (
	"gorm-gen-skeleton/app/amqp/consumer"
	"gorm-gen-skeleton/internal/mq"
)

type Amqp struct{}

func (*Amqp) InitConsumers() []mq.ConsumerHandler {
	return []mq.ConsumerHandler{
		&consumer.FooConsumer{},
	}
}
