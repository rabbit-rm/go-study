package rabbitmq

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestConsumerPriority(t *testing.T) {
	consumerPriorityServer()
	go consumerPriorityClient("consumer5", 10)
	go consumerPriorityClient("consumer1", 10)
	go consumerPriorityClient("consumer2", 5)
	go consumerPriorityClient("consumer3", 1)
	select {}
}

func consumerPriorityServer() {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}

	queue, err := ch.QueueDeclare("priority_queue", false, false, false, false, nil)
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
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func consumerPriorityClient(consumerTag string, priority int) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}
	messages, err := ch.Consume("priority_queue", consumerTag, true, false, false, false, amqp091.Table{
		"x-priority": priority,
	})
	if err != nil {
		log.Fatalf("consume queue error:%+v\n", err)
	}

	n := rand.Int31n(50) * int32(priority)

	for msg := range messages {
		if n <= 0 {
			log.Printf("[%s] break\n", consumerTag)
			break
		}
		log.Printf("[%s] receive:%s\n", consumerTag, msg.Body)
		n--
	}

	_ = ch.Cancel(consumerTag, true)
}
