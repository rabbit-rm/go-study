package rabbitmq

import (
	"fmt"
	"log"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

func TestRPCDirectTo(t *testing.T) {
	go func() {
		err := rpcReplyToServer()
		if err != nil {
			log.Fatalf("error:%+v\n", err)
		}
	}()

	_ = rpcReplyToClient()
	select {}
}

func rpcReplyToServer() error {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	queue, err := ch.QueueDeclare("rpc", false, false, false, false, nil)
	if err != nil {
		return err
	}
	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	for msg := range messages {
		log.Printf("receive msg:%s,%s\n", msg.Body, msg.ReplyTo)
		_ = ch.Publish("", msg.ReplyTo, false, false, amqp091.Publishing{
			ContentType:     msg.ContentType,
			ContentEncoding: msg.ContentEncoding,
			Body:            msg.Body,
		})
	}
	return nil
}

func rpcReplyToClient() error {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	queue, err := ch.QueueDeclare("result", false, false, false, false, nil)
	if err != nil {
		return err
	}
	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		_ = ch.Publish("", "rpc", false, false, amqp091.Publishing{
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			ReplyTo:         queue.Name,
			Body:            []byte(fmt.Sprintf("hello[%d]", i)),
		})
	}
	for msg := range messages {
		log.Printf("[]receive msg:%s\n", msg.Body)
	}
	return nil
}
