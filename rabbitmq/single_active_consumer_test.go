package rabbitmq

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestSingleActive(t *testing.T) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create conn err:%+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel err:%+v\n", err)
	}

	defer func() {
		_ = ch.Close()
	}()

	_, err = ch.QueueDelete("single_active_consumer_queue", false, false, false)
	if err != nil {
		log.Fatalf("create queue err:%+v\n", err)
	}

	queue, err := ch.QueueDeclare("single_active_consumer_queue", false, false, false, false, amqp091.Table{
		"x-single-active-consumer": true,
	})
	if err != nil {
		log.Fatalf("declare queue err:%+v\n", err)
	}

	go singleActiveConsumerPublisher(queue, ch)

	go declareSingleActiveConsumer("consumer-1")
	go declareSingleActiveConsumer("consumer-2")
	go declareSingleActiveConsumer("consumer-3")
	go declareSingleActiveConsumer("consumer-4")
	go declareSingleActiveConsumer("consumer-5")
	go declareSingleActiveConsumer("consumer-6")
	go declareSingleActiveConsumer("consumer-7")

	select {}

}

func singleActiveConsumerPublisher(queue amqp091.Queue, ch *amqp091.Channel) {
	index := 0
	for {
		_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			Body:            []byte(fmt.Sprintf("hello[%d]", index)),
		})
		time.Sleep(100 * time.Millisecond)
		index++
	}
}

func declareSingleActiveConsumer(consumerTag string) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create conn err:%+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel err:%+v\n", err)
	}

	defer func() {
		_ = ch.Close()
	}()

	messages, err := ch.Consume("single_active_consumer_queue", consumerTag, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("create consumer err:%+v\n", err)
	}

	n := rand.Int31n(50)
	nn := n
	for msg := range messages {
		log.Printf("[%s]receive msg:%s\n", msg.ConsumerTag, msg.Body)
		_ = ch.Ack(msg.DeliveryTag, false)
		n--
		if n <= 0 {
			log.Printf("%s[%d] broken", msg.ConsumerTag, nn)
			break
		}
	}
	err = ch.Cancel(consumerTag, true)
	if err != nil {
		log.Fatalf("err:%+v\n", err)
	}
}
