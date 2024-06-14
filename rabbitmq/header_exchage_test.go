package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestHeaderExchange(t *testing.T) {
	go headerExchangePublisher("key1", "147")
	// 消息丢失
	go headerExchangePublisher("key1", "258")
	// 消息丢失
	go headerExchangePublisher("key2", "258")
	go headerExchangePublisher("key3", "369")

	go headerExchangeConsumer("key1", "147")
	go headerExchangeConsumer("key3", "369")
	go headerExchangeConsumer("x-match", "any")
	select {}
}

func headerExchangePublisher(key, value string) {
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
	err = ch.ExchangeDeclare("header_exchange", amqp091.ExchangeHeaders, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare exchange error: %+v\n", err)
	}
	/*queue, err := ch.QueueDeclare("header_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}
	var table = map[string]interface{}{
		"key1":    "147",
		"key2":    "258",
		"key3":    "369",
		"x-match": "all",
	}
	err = ch.QueueBind(queue.Name, "", "header_exchange", false, table)
	if err != nil {
		log.Fatalf("bind queue error: %+v\n", err)
	}*/
	index := 1
	for {
		err = ch.Publish("header_exchange", "", false, false, amqp091.Publishing{
			Headers: map[string]interface{}{
				key: value,
			},
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("header message %d", index)),
		})
		if err != nil {
			log.Printf("publish error: %+v\n", err)
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func headerExchangeConsumer(key, value string) {
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
	if err := ch.QueueBind(queue.Name, "", "header_exchange", false, map[string]interface{}{
		key: value,
	}); err != nil {
		log.Fatalf("bind queue error: %+v\n", err)
	}
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %+v\n", err)
	}
	for msg := range msgs {
		log.Printf("receive message: [%s:%s]header:%s,body:%s\n", key, value, msg.Headers, msg.Body)
	}
}
