package rabbitmq

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestWorkQueue2(t *testing.T) {
	queueName := "work_queue_persistence"
	go workQueueProducer2("producer_01", queueName)
	// go workQueueProducer("producer_02", queueName)
	go workQueueConsumer2("consumer_01", queueName)
	go workQueueConsumer2("consumer_02", queueName)
	go workQueueConsumer2("consumer_03", queueName)
	select {}
}

func workQueueProducer2(producerName, queueName string) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}
	// 定义持久化队列
	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("create work queue error:%+v\n", err)
	}
	messageGen := func(separator string, num int) string {
		var b bytes.Buffer
		for range num {
			b.Write([]byte(separator))
		}
		return b.String()
	}
	for i := 0; i < 100; i++ {
		messageId := fmt.Sprintf("%s_%d", producerName, i)
		err := channel.Publish(amqp091.DefaultExchange, queue.Name, false, false, amqp091.Publishing{
			// 定义消息为持久消息
			DeliveryMode:    amqp091.Persistent,
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			MessageId:       messageId,
			Timestamp:       time.Now(),
			Body:            []byte(messageGen(".", i%10+1)),
		})
		if err != nil {
			log.Printf("publish message error:%+v\n", err)
		}
	}
}

func workQueueConsumer2(consumerName, queueName string) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}
	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("create work queue error:%+v\n", err)
	}
	// 设置预计取数为1，一次只向消费者发送一个消息，并等待消费者处理并确认消息之后再次发送下一条消息
	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("qos errro:%+v\n", err)
	}
	consume, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("create consume error:%+v\n", err)
	}
	for message := range consume {
		id := message.MessageId
		body := message.Body
		count := bytes.Count(body, []byte("."))
		time.Sleep(time.Duration(count) * time.Millisecond * 500)
		log.Printf("Receive Message: %s,[%s]%s\n", consumerName, id, body)
		_ = message.Ack(false)
	}
}

func TestWorkQueue(t *testing.T) {
	queueName := "work_queue"
	go workQueueProducer("producer_01", queueName)
	// go workQueueProducer("producer_02", queueName)
	go workQueueConsumer("consumer_01", queueName)
	go workQueueConsumer("consumer_02", queueName)
	go workQueueConsumer("consumer_03", queueName)
	select {}
}

func workQueueProducer(producerName, queueName string) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}
	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("create work queue error:%+v\n", err)
	}
	messageGen := func(separator string, num int) string {
		var b bytes.Buffer
		for range num {
			b.Write([]byte(separator))
		}
		return b.String()
	}
	for i := 0; i < 10; i++ {
		messageId := fmt.Sprintf("%s_%d", producerName, i)
		err := channel.Publish(amqp091.DefaultExchange, queue.Name, false, false, amqp091.Publishing{
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			MessageId:       messageId,
			Timestamp:       time.Now(),
			Body:            []byte(messageGen(".", i%5+1)),
		})
		if err != nil {
			log.Printf("publish message error:%+v\n", err)
		}
	}
}

func workQueueConsumer(consumerName, queueName string) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}
	_, err = channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("create work queue error:%+v\n", err)
	}
	consume, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("create consume error:%+v\n", err)
	}
	for message := range consume {
		id := message.MessageId
		body := message.Body
		count := bytes.Count(body, []byte("."))
		time.Sleep(time.Duration(count) * time.Millisecond * 500)
		log.Printf("Receive Message: %s,[%s]%s\n", consumerName, id, body)
		_ = message.Ack(false)
	}
}
