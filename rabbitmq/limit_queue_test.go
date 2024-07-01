package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestLimitQueueDefine(t *testing.T) {
	conn := limitQueueDefine()

	go func() {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}
		index := 1
		for {
			_ = ch.Publish("limit_task_exchange", "limit-task", false, false, amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("limit task message %d", index)),
			})
			time.Sleep(1 * time.Millisecond)
			index++
		}
	}()

	go func() {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}
		err = ch.Qos(10, 0, false)
		if err != nil {
			log.Fatalf("qos error: %+v", err)
		}
		messages, err := ch.Consume("limit_task_queue", "", false, false, false, false, nil)
		for message := range messages {
			log.Printf("A consumer message: %+v", string(message.Body))
			_ = message.Ack(false)
		}
	}()

	go func() {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}
		err = ch.Qos(20, 0, false)
		if err != nil {
			log.Fatalf("qos error: %+v", err)
		}
		messages, err := ch.Consume("limit_task_queue", "", false, false, false, false, nil)
		for message := range messages {
			log.Printf("B consumer message: %+v", string(message.Body))
			_ = message.Ack(false)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}
	err = ch.Qos(50, 0, false)
	if err != nil {
		log.Fatalf("qos error: %+v", err)
	}
	messages, err := ch.Consume("limit_task_queue", "", false, false, false, false, nil)
	for message := range messages {
		log.Printf("C consumer message: %+v", string(message.Body))
		_ = message.Ack(false)
	}

}

func limitQueueDefine() *amqp091.Connection {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}

	ch, err := conn.Channel()
	_, err = ch.QueueDelete("limit_task_queue", false, false, false)
	if err != nil {
		log.Fatalf("queue delete error: %+v", err)
	}
	queue, err := ch.QueueDeclare("limit_task_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare error: %+v", err)
	}

	err = ch.ExchangeDelete("limit_task_exchange", false, false)
	if err != nil {
		log.Fatalf("exchange delete error: %+v", err)
	}
	err = ch.ExchangeDeclare("limit_task_exchange", amqp091.ExchangeDirect, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange declare error: %+v", err)
	}

	err = ch.QueueBind(queue.Name, "limit-task", "limit_task_exchange", false, nil)
	if err != nil {
		log.Fatalf("queue bind error: %+v", err)
	}

	defer func() {
		_ = ch.Close()
	}()

	return conn
}
