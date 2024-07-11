package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
)

func connection() (*amqp091.Connection, error) {
	return amqp091.Dial(rabbitUrl)
}
