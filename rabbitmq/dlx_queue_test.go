package rabbitmq

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestDlxQueue(t *testing.T) {
	conn := dlxQueueDefine()
	go func() {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}
		index := 1
		for {
			err := ch.Publish("my_normal_exchange", "test.normal.1", false, false, amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("hello world %d", index)),
				Expiration:  strconv.Itoa(index * 10),
			})
			if err != nil {
				log.Printf("publish message error: %+v", err)
			}
			index++
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}
		ch.Qos(1, 0, false)
		messages, err := ch.Consume("my_normal_queue", "", false, false, false, false, nil)
		if err != nil {
			log.Fatalf("consume message error: %+v", err)
		}

		for message := range messages {
			log.Printf("consume message: %s", message.Body)
			time.Sleep(600 * time.Millisecond)
			_ = message.Nack(true, false)
		}
	}()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}
	ch.Qos(1, 0, false)
	messages, err := ch.Consume("dlx_queue", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume message error: %+v", err)
	}

	for message := range messages {
		log.Printf("consume dlx queue message: %s", message.Body)
	}

}

func dlxQueueDefine() *amqp091.Connection {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}
	err = ch.ExchangeDeclare("my_normal_exchange", amqp091.ExchangeTopic, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare exchange error: %+v", err)
	}

	_, err = ch.QueueDelete("my_normal_queue", false, false, false)
	if err != nil {
		log.Fatalf("delete queue error: %+v", err)
	}

	queue, err := ch.QueueDeclare("my_normal_queue", false, false, false, false, amqp091.Table{
		"x-dead-letter-exchange":    "dlx_exchange",
		"x-dead-letter-routing-key": "dlx.message",
		"x-message-ttl":             500,
		"x-max-length":              20,
	})
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}

	err = ch.QueueBind(queue.Name, "test.normal.*", "my_normal_exchange", false, nil)
	if err != nil {
		log.Fatalf("bind queue error: %+v", err)
	}

	// 定义死信交换机
	_ = ch.ExchangeDelete("dlx_exchange", false, false)
	err = ch.ExchangeDeclare("dlx_exchange", amqp091.ExchangeTopic, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare exchange error: %+v", err)
	}
	// 定义死信队列
	_, err = ch.QueueDelete("dlx_queue", false, false, false)
	if err != nil {
		log.Fatalf("delete queue error: %+v", err)
	}
	dlxQueue, err := ch.QueueDeclare("dlx_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare dlx queue error: %+v", err)
	}

	err = ch.QueueBind(dlxQueue.Name, "dlx.*", "dlx_exchange", false, nil)
	if err != nil {
		log.Fatalf("bind queue error: %+v", err)
	}

	defer func() {
		_ = ch.Close()
	}()

	return conn

}
