package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestMessagePriority(t *testing.T) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connetion error: %+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// server
	{
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v\n", err)
		}

		defer func() {
			_ = ch.Close()
		}()

		queue, err := ch.QueueDeclare("priority_queue", false, false, true, false, nil)
		if err != nil {
			log.Fatalf("declare queue error: %+v\n", err)
		}

		go func() {
			index := 1
			for {
				_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
					Priority:    uint8(10 - index%10),
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("hello[%d]", index)),
				})
				index++
				time.Sleep(200 * time.Millisecond)
			}
		}()
	}

	// client
	{
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v\n", err)
		}

		messages, err := ch.Consume("priority_queue", "", false, false, false, false, nil)
		if err != nil {
			log.Fatalf("consume queue error: %+v\n", err)
		}

		for msg := range messages {
			log.Printf("receive message: [%d]%s\n", msg.Priority, msg.Body)
		}
	}

}
