package rabbitmq

import (
	"log"
	"strconv"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

func TestRPC(t *testing.T) {
	go rpcServer()
	go rpcClient(11)

	select {}

}

func rpcServer() {
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
	queue, err := ch.QueueDeclare("rpc_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("set qos error: %+v\n", err)
	}
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %+v\n", err)
	}
	for msg := range msgs {
		log.Printf("receive message: %s\n", msg.Body)
		n, err := strconv.Atoi(string(msg.Body))
		if err != nil {
			log.Fatalf("convert message error: %+v\n", err)
		}
		err = ch.Publish("", msg.ReplyTo, false, false, amqp091.Publishing{
			ContentType:   "text/plain",
			CorrelationId: msg.CorrelationId,
			Body:          []byte(strconv.Itoa(fib(n))),
		})
		if err != nil {
			log.Printf("publish error: %+v\n", err)
		}
	}
}

func rpcClient(num int) {
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
	queue, err := ch.QueueDeclare("", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %+v\n", err)
	}
	err = ch.Publish("", "rpc_queue", false, false, amqp091.Publishing{
		ContentType:   "text/plain",
		CorrelationId: "1",
		ReplyTo:       queue.Name,
		Body:          []byte(strconv.Itoa(num)),
	})
	if err != nil {
		log.Fatalf("publish error: %+v\n", err)
	}
	for msg := range msgs {
		if msg.CorrelationId == "1" {
			resp, err := strconv.Atoi(string(msg.Body))
			if err != nil {
				log.Printf("convert response error: %+v\n", err)
			}
			log.Printf("receive message1: %d\n", resp)
			break
		}
	}
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}
