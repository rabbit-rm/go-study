package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestStreamsSingleActiveConsumerServer(t *testing.T) {
	conn, err := connection()
	if err != nil {
		log.Fatalf("create conn error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// declare streams
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	queue, err := ch.QueueDeclare(
		"streams_single_queue",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type": "stream",
			// not support
			"x-single-active-consumer": true,
		})
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}

	index := 1
	for {

		_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("message-%d", index)),
		})

		time.Sleep(200 * time.Millisecond)
		index++
	}
}

func TestStreamSingleClient1(t *testing.T) {
	streamSingleClient("stream-consumer-1")
}

func TestStreamSingleClient2(t *testing.T) {
	streamSingleClient("stream-consumer-2")
}

func TestStreamSingleClient3(t *testing.T) {
	streamSingleClient("stream-consumer-3")
}
func TestStreamSingleClient4(t *testing.T) {
	streamSingleClient("stream-consumer-4")
}

func streamSingleClient(consumerName string) {
	conn, err := connection()
	if err != nil {
		log.Fatalf("create conn error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	// consumer
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	_ = ch.Qos(100, 0, false)

	messages, err := ch.Consume(
		"streams_single_queue",
		consumerName,
		false,
		false,
		false,
		false,
		amqp091.Table{
			// 从接入后第一条消息开始消费
			"x-stream-offset": "next",
		})
	if err != nil {
		log.Fatalf("consume error: %+v", err)
	}

	for msg := range messages {
		log.Printf("consumer message: %s,timestamp:%s", msg.Body, time.Now().Sub(msg.Timestamp))
		_ = msg.Ack(false)
	}
}
