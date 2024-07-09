package rabbitmq

import (
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestQueueMessageTTL(t *testing.T) {

	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.132:5672")
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	// 定义消息ttl 队列
	queue, err := ch.QueueDeclare("message_ttl_queue", false, true, true, false, amqp091.Table{
		"x-message-ttl": 5000,
	})
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}

	go func() {

		for {
			_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
				ContentType: "text/plain",
				// 发布消息 设置消息 ttl
				Expiration: "3000",
				Body:       []byte("hello world"),
			})
			time.Sleep(200 * time.Millisecond)
		}
	}()

	messages, err := ch.Consume(queue.Name, "", true, true, false, false, nil)
	if err != nil {
		log.Fatalf("consume message error: %+v", err)
	}

	for msg := range messages {
		log.Printf("consume message: body:%s,expiration:%+v", string(msg.Body), msg.Expiration)
	}
}

func TestTTLQueue(t *testing.T) {

	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.132:5672")
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	// 定义 ttl队列，没有消费者的情况下，队列过期时间（此处 5s）到了，队列会被删除
	queue, err := ch.QueueDeclare("ttl_queue", false, false, false, false, amqp091.Table{
		"x-expires": 5000,
	})
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}

	go func() {

		for {
			_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
				ContentType: "text/plain",
				// 发布消息 设置消息 ttl
				Expiration: "3000",
				Body:       []byte("hello world"),
			})
			time.Sleep(200 * time.Millisecond)
		}
	}()

	messages, err := ch.Consume(queue.Name, "", true, true, false, false, nil)
	if err != nil {
		log.Fatalf("consume message error: %+v", err)
	}
	index := 1
	consumerName := ""
	for msg := range messages {
		if strings.Compare(consumerName, "") == 0 {
			consumerName = msg.ConsumerTag
		}
		log.Printf("consume message: body:%s,expiration:%+v", string(msg.Body), msg.Expiration)
		index++
		if index > 100 {
			break
		}
	}
	_ = ch.Cancel(consumerName, true)

	time.Sleep(5 * time.Second)
}

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
