package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestHelloWorld(t *testing.T) {
	go producers()
	consumers()
}

func producers() {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("conn error:%+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	// queue name
	const queueName = "rabbit.rm"
	var wg sync.WaitGroup
	// 10 producer
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(conn *amqp091.Connection) {
			defer wg.Done()
			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("create channel error:%+v\n", err)
			}
			_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
			if err != nil {
				log.Fatalf("create queue error:%+v\n", err)
			}
			for i := 0; i < 100000; i++ {
				err := ch.Publish("", queueName, false, false, amqp091.Publishing{
					ContentType:     "text/plain",
					ContentEncoding: "utf-8",
					Timestamp:       time.Now(),
					Body:            []byte(fmt.Sprintf("[%d]hello world!", i)),
				})
				if err != nil {
					log.Printf("send msg error:%v\n", err)
					continue
				}
				// time.Sleep(time.Duration(rand.Int31n(10)*5) * time.Millisecond)
			}
		}(conn)
	}
	wg.Wait()
}

func consumers() {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("conn error:%+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	// queue name
	const queueName = "rabbit.rm"
	var wg sync.WaitGroup
	// 10 producer
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(conn *amqp091.Connection) {
			defer wg.Done()
			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("create channel error:%+v\n", err)
			}
			_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
			if err != nil {
				log.Fatalf("create queue error:%+v\n", err)
			}
			consume, err := ch.Consume(queueName, "", true, false, false, false, nil)
			if err != nil {
				log.Fatalf("create consumer error:%+v\n", err)
			}
			for msg := range consume {
				log.Printf("receive msg:%s\n", msg.Body)
			}
		}(conn)
		wg.Wait()
	}
}

func TestSimpleHelloWorld(t *testing.T) {
	go simpleProducer()
	simpleConsumer()
}

func simpleProducer() {
	// 连接 rabbitmq 服务
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("dial error:%+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	// 获得一个 channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel error:%+v", err)
	}
	// 定义一个队列
	queue, err := ch.QueueDeclare("hello-world", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error:%+v", err)
	}
	// 向队列发送消息
	for i := 0; i < 100000; i++ {
		err = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
			Headers:     nil,
			ContentType: "text/plain",
			Body:        ([]byte)("hello world!"),
		})
		if err != nil {
			log.Fatalf("publish message error:%+v", err)
		}
	}
}

func simpleConsumer() {
	// 创建 rabbitmq 连接
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("dial error:%+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	// 创建一个 channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel error:%+v", err)
	}
	// 消费者一般不需要定义队列，防止生产者在消费者之后启动，这里定义同样队列
	queue, err := ch.QueueDeclare("hello-world", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error:%+v", err)
	}
	// 消费 && 处理 消息
	consume, err := ch.Consume(queue.Name, "hello-world", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("simpleConsumer message error:%+v", err)
	}
	for msg := range consume {
		log.Printf("receive msg:%v\n", msg.Body)
	}
}
