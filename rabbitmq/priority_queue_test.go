package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestQueuePriority(t *testing.T) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	{
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}

		queue, err := ch.QueueDeclare(
			"priority-queue",
			false,
			true,
			true,
			false,
			amqp091.Table{
				"x-max-priority": 5,
			})
		if err != nil {
			log.Fatalf("declare queue error: %+v", err)
		}

		go func() {
			index := 1
			for {
				_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("message-%d", index)),
					Priority:    uint8(index % 10),
				})
				time.Sleep(100 * time.Millisecond)
				index++
			}
		}()
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	_ = ch.Qos(10, 0, false)
	messages, err := ch.Consume("priority-queue", "priority-consumer", false, true, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %+v", err)
	}
	for msg := range messages {
		log.Printf("consume message: priority:%d,body:%s", msg.Priority, msg.Body)
		time.Sleep(150 * time.Millisecond)
		_ = msg.Ack(false)
	}

}
