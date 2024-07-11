package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestStreamsServer(t *testing.T) {
	conn, err := connection()
	if err != nil {
		log.Fatalf("create conn error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// declare streams
	var queueName string
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	queue, err := ch.QueueDeclare(
		"streams_queue_1",
		true,
		false,
		false,
		false,
		amqp091.Table{
			// 定义队列类型为 stream
			"x-queue-type": "stream",
			// 定义 Streams 的最大长度
			"x-max-length-bytes": 1024,
			// 定义 Streams 消息最大生命周期 1min
			"x-max-age": "10s",
		})
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}
	queueName = queue.Name

	index := 1
	for {
		_ = ch.Publish("", queueName, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("message-%d", index)),
			Timestamp:   time.Now(),
		})
		time.Sleep(10 * time.Millisecond)
		index++
	}

}

func TestStreamsClient1(t *testing.T) {
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

	// 1 min 之前
	// timestamp := time.Now().Add(-1 * time.Minute).UnixMilli()

	messages, err := ch.Consume(
		"streams_queue_1",
		"streams_consumer_1",
		false,
		false,
		false,
		false,
		amqp091.Table{
			// 从第一条消息开始消费
			// "x-stream-offset": "first",
			// 从接入后第一条消息开始消费
			// "x-stream-offset": "last",
			// 从第 20000 个消息开始消费
			// "x-stream-offset": 20000,
			// 从 5s 前的消息开始消费
			// "x-stream-offset": "5s",
			// 从指定时间戳的时间的消息开始消费
			"x-stream-offset": time.Now().Add(-1 * time.Minute),
		})
	if err != nil {
		log.Fatalf("consume error: %+v", err)
	}

	for msg := range messages {
		log.Printf("consumer message: %s,timestamp:%s", msg.Body, time.Now().Sub(msg.Timestamp))
		_ = msg.Ack(false)
	}
}
