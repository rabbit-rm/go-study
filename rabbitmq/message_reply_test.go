package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestMessageReply(t *testing.T) {
	go messageReplyPublisher()
	go messageReplyConsumer("A", true)
	go messageReplyConsumer("B", false)
	go messageReplyConsumer("C", false)
	select {}
}

func messageReplyPublisher() {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673")
	if err != nil {
		log.Fatalf("create connection error: %+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v\n", err)
	}
	queue, err := ch.QueueDeclare("message_reply_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}
	index := 1
	for {
		_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("test message %d", index)),
		})
		index++
		time.Sleep(200 * time.Millisecond)
	}
}

func messageReplyConsumer(name string, isReply bool) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673")
	if err != nil {
		log.Fatalf("create connection error: %+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v\n", err)
	}
	queue, err := ch.QueueDeclare("message_reply_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}
	messages, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %+v\n", err)
	}
	for msg := range messages {
		if isReply {
			log.Printf("receive message:[%s] %s\n", name, msg.Body)
			// multiple 表示是否批量应答
			_ = msg.Ack(false)
		} else {
			// 批量拒绝，requeue 表示是否重新入队
			// _ = msg.Nack(false, false)
			// requeue 表示是否重新入队
			log.Printf("reject message:[%s] %s\n", name, msg.Body)
			_ = msg.Reject(true)
		}
	}
}
