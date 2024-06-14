package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestMessageQos(t *testing.T) {
	go messageQosPublisher()
	go messageQosConsumer("A", time.Second, 5, 0)
	go messageQosConsumer("B", 5*time.Second, 3, 0)
	// go messageQosConsumer("C", 10*time.Second, 2, 0)
	select {}
}

func messageQosPublisher() {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("create connection error: %+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v\n", err)
	}
	queue, err := ch.QueueDeclare("work_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}
	for i := 0; i < 10; i++ {
		_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("test message %d", i)),
		})
	}
}

func messageQosConsumer(name string, duration time.Duration, prefetchCount, prefetchSize int) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("create connection error: %+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v\n", err)
	}
	err = ch.Qos(prefetchCount, prefetchSize, false)
	if err != nil {
		log.Fatalf("set qos error: %+v\n", err)
	}
	queue, err := ch.QueueDeclare("work_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}
	messages, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %+v\n", err)
	}
	for msg := range messages {
		time.Sleep(duration)
		log.Printf("receive message: [%s]%s\n", name, msg.Body)
		_ = msg.Ack(false)
	}
}
