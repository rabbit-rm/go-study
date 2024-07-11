package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestStreamsFilter(t *testing.T) {
	conn, err := connection()
	if err != nil {
		log.Fatalf("create conn error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	var queueName string
	{
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}

		queue, err := ch.QueueDeclare(
			"streams_filter_queue",
			true,
			false,
			false,
			false,
			amqp091.Table{
				"x-queue-type": "stream",
			})
		if err != nil {
			log.Fatalf("declare queue error: %+v", err)
		}
		queueName = queue.Name

		go func() {
			index := 1
			for {
				_ = ch.Publish("", queueName, false, false, amqp091.Publishing{
					Headers: amqp091.Table{
						"x-stream-filter-value": fmt.Sprintf("filter-%d", index%5+1),
					},
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("message-%d", index)),
				})
				time.Sleep(20 * time.Millisecond)
				index++
			}
		}()
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	_ = ch.Qos(100, 0, false)

	messages, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		amqp091.Table{
			// "x-stream-filter":           "filter-4",
			"x-stream-match-unfiltered": true,
		})
	if err != nil {
		log.Fatalf("consume error: %+v", err)
	}

	for msg := range messages {
		// 因为服务器端是概率性的，所以客户端也需要过滤
		if msg.Headers["x-stream-filter-value"] == "filter-4" {
			log.Printf("consume message: body:%s", msg.Body)
			_ = msg.Ack(false)
		}
	}
}
