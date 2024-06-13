package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestLogSystem(t *testing.T) {
	go LogProducer()
	go LogConsumer(func(msg []byte) {
		log.Printf("receive message:%s\n", msg)
	})
	out, err := os.OpenFile("W:\\GoProject\\private\\study\\rabbitmq\\logs.out", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("create out error:%+v\n", err)
	}
	go LogConsumer(func(msg []byte) {
		_, _ = out.Write(msg)
		_, _ = out.Write([]byte("\n"))
	})
	select {}
}

func LogProducer() {
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
	// 定义一个交换机
	// 类型是：fanout（amqp091.ExchangeFanout），fanout 仅仅只适用于广播消息
	if err := channel.ExchangeDeclare("logs", amqp091.ExchangeFanout, true, false, false, false, nil); err != nil {
		log.Fatalf("declare exchange error:%+v\n", err)
	}
	index := 1
	for {
		msg := fmt.Sprintf("message_%2d", index)
		if err = channel.Publish("logs", "", false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		}); err != nil {
			log.Printf("publish message error:%+v\n", err)
			continue
		}
		log.Printf("publish message:%s\n", msg)
		index++
		time.Sleep(time.Millisecond * 100)
	}
}

func LogConsumer(f func(msg []byte)) {
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
	if err := channel.ExchangeDeclare("logs", amqp091.ExchangeFanout, true, false, false, false, nil); err != nil {
		log.Fatalf("declare exchange error:%+v\n", err)
	}
	// 定一个独占队列
	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("declare queue error:%+v\n", err)
	}
	err = channel.QueueBind(queue.Name, "", "logs", false, nil)
	if err != nil {
		log.Fatalf("bind queue error:%+v\n", err)
	}
	consume, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("create consumer error:%+v\n", err)
	}
	forever := make(chan struct{})
	go func() {
		for msg := range consume {
			f(msg.Body)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
