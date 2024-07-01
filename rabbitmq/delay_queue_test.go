package rabbitmq

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestDelayQueue(t *testing.T) {
	conn := delayQueueDefine()

	go func() {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}
		index := 1
		for {
			_ = ch.Publish("", "task_queue", false, false, amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("task message %d", index)),
				Timestamp:   time.Now(),
				Expiration:  strconv.Itoa(int(rand.Int31n(5) * 1000)),
			})
			time.Sleep(100 * time.Millisecond)
			index++
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}
	messages, err := ch.Consume("delay_queue", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %+v", err)
	}

	for message := range messages {
		log.Printf("delay consumer message: [%s]%+v,expiration:%sms",
			message.Timestamp.Format("2006-01-02:15:04:05.000"),
			string(message.Body),
			message.Expiration)
		_ = message.Ack(false)
	}

}

func delayQueueDefine() *amqp091.Connection {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}

	// 正常队列
	ch, err := conn.Channel()
	_, err = ch.QueueDelete("task_queue", false, false, false)
	if err != nil {
		log.Fatalf("queue delete error: %+v", err)
	}
	_, err = ch.QueueDeclare("task_queue", false, false, false, false, amqp091.Table{
		"x-dead-letter-exchange":    "delay_exchange",
		"x-dead-letter-routing-key": "task.delay.5s",
	})
	if err != nil {
		log.Fatalf("queue declare error: %+v", err)
	}

	// 延迟队列
	_, err = ch.QueueDelete("delay_queue", false, false, false)
	if err != nil {
		log.Fatalf("queue delete error: %+v", err)
	}
	delayQueue, err := ch.QueueDeclare("delay_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare error: %+v", err)
	}
	err = ch.ExchangeDelete("delay_exchange", false, false)
	if err != nil {
		log.Fatalf("exchange delete error: %+v", err)
	}
	err = ch.ExchangeDeclare("delay_exchange", amqp091.ExchangeTopic, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange declare error: %+v", err)
	}

	err = ch.QueueBind(delayQueue.Name, "task.delay.*", "delay_exchange", false, nil)
	if err != nil {
		log.Fatalf("queue bind error: %+v", err)
	}

	defer func() {
		_ = ch.Close()
	}()

	return conn
}
