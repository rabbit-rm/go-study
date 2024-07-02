package rabbitmq

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestCancelConsumer(t *testing.T) {
	consumerCancelServer()
	consumerCancelClient()
}

func consumerCancelServer() {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}

	queue, err := ch.QueueDeclare("cancel_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error:%+v\n", err)
	}

	go func() {
		index := 1
		for {
			_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("hello[%d]", index)),
			})
			index++
			time.Sleep(100 + time.Millisecond)
		}
	}()
}

func consumerCancelClient() {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}
	cancelCh := make(chan string)
	ch.NotifyCancel(cancelCh)
	go func() {
		for c := range cancelCh {
			log.Printf("consumer cancel:%s\n", c)
		}
	}()

	messages, err := ch.Consume("cancel_queue", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume queue error:%+v\n", err)
	}

	n := rand.Int31n(30)
	for msg := range messages {
		if n <= 0 {
			log.Printf("consumer break\n")
			break
		}
		log.Printf("receive message:%s\n", string(msg.Body))
		n--
	}

	_, _ = ch.QueueDelete("cancel_queue", false, false, false)
	select {}

}
