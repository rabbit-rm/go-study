package rabbitmq

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestQueueTTL(t *testing.T) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	publisherCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}
	// 队列设置过期时间
	queue, err := publisherCh.QueueDeclare("queue_ttl", false, false, false, false, amqp091.Table{
		"x-message-ttl": 500,
	})
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}
	for i := 0; i < 20; i++ {
		err = publisherCh.Publish("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte("hello world"),
			Expiration:  strconv.Itoa((20 - i) * 100),
		})
		if err != nil {
			log.Fatalf("publish message error: %+v", err)
		}
	}

	//
	time.Sleep(500 * time.Millisecond)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume message error: %+v", err)
	}
	for message := range messages {
		log.Printf("consume message: contentType:%+v,body:%s,expiration:%+v", message.ContentType, string(message.Body), message.Expiration)
	}
}
