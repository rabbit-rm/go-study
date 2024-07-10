package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestQueueLengthLimitedOverflow(t *testing.T) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.132:5672")
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// declare dlx
	var dlxName string
	{
		dlxName = "queue_limited_overflow_dlx"
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}
		err = ch.ExchangeDeclare(dlxName, "direct", false, false, false, false, nil)
		if err != nil {
			log.Fatalf("declare dlx error: %+v", err)
		}

		queue, err := ch.QueueDeclare("overflow_dlx_queue", false, false, false, false, nil)
		if err != nil {
			log.Fatalf("declare dlx queue error: %+v", err)
		}

		err = ch.QueueBind(queue.Name, "overflow_dlx", dlxName, false, nil)
		if err != nil {
			log.Fatalf("bind dlx queue error: %+v", err)
		}

		go func() {
			messages, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
			if err != nil {
				log.Fatalf("consume dlx message error: %+v", err)
			}

			for msg := range messages {
				log.Printf("consume dlx message: body:%s", string(msg.Body))
			}
		}()
	}
	var queueName string

	{
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}

		queue, err := ch.QueueDeclare("queue_length_limited_overflow", false, false, false, false, amqp091.Table{
			"x-max-length":       10,
			"x-max-length-bytes": 1024 * 1024,
			// 溢出返回给发布者，消息不丢失
			// "x-overflow":             "reject-publish",
			// 溢出消息死信，消息发送到死信队列
			"x-overflow":                "reject-publish-dlx",
			"x-dead-letter-exchange":    dlxName,
			"x-dead-letter-routing-key": "overflow_dlx",
		})
		if err != nil {
			log.Fatalf("declare queue error: %+v", err)
		}
		queueName = queue.Name
		go func() {
			index := 1
			for {
				_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("hello[%d]", index)),
				})
				time.Sleep(200 * time.Millisecond)
				index++
			}
		}()
	}
	time.Sleep(15 * 200 * time.Millisecond)
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	messages, err := ch.Consume(queueName, "", false, true, false, false, nil)
	if err != nil {
		log.Fatalf("consume message error: %+v", err)
	}

	for msg := range messages {
		log.Printf("consume message: body:%s", string(msg.Body))
	}
}

func TestQueueLengthLimited(t *testing.T) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.132:5672")
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
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

		queue, err := ch.QueueDeclare("queue_length_limited", false, false, false, false, amqp091.Table{
			"x-max-length":       10,
			"x-max-length-bytes": 1024 * 1024,
		})
		if err != nil {
			log.Fatalf("declare queue error: %+v", err)
		}
		queueName = queue.Name
		go func() {
			index := 1
			for {
				_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("hello[%d]", index)),
				})
				time.Sleep(200 * time.Millisecond)
				index++
			}
		}()
	}
	time.Sleep(15 * 200 * time.Millisecond)
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}

	messages, err := ch.Consume(queueName, "", false, true, false, false, nil)
	if err != nil {
		log.Fatalf("consume message error: %+v", err)
	}

	for msg := range messages {
		log.Printf("consume message: body:%s", string(msg.Body))
	}
}
