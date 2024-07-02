package rabbitmq

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestExclusiveQueue(t *testing.T) {
	exclusiveQueueServer()
}

func exclusiveQueueServer() {
	const consumerTag = "exclusive-consumer-1"
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}

	queue, err := ch.QueueDeclare("exclusive_queue", false, false, true, false, nil)
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

	ch2, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}

	messages, err := ch2.Consume(queue.Name, consumerTag, true, true, false, false, nil)
	if err != nil {
		log.Fatalf("consume queue error:%+v\n", err)
	}

	n := rand.Int31n(100)
	for msg := range messages {
		if n <= 0 {
			log.Printf("%s break\n", consumerTag)
			break
		}
		log.Printf("%s receive:%s\n", consumerTag, msg.Body)
		n--
	}

	_ = ch.Cancel(consumerTag, false)

}
