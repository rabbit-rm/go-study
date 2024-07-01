package rabbitmq

import (
	"fmt"
	"log"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

func TestAlternativeExchange(t *testing.T) {
	go alternativePublisher()

	go alternativeConsumer("route_a")
	go alternativeConsumer("route_e")

	go alternativeExchangeConsumer("route_c")

	select {}

}

func alternativePublisher() {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}
	err = ch.ExchangeDeclare("my_alternative_exchange", amqp091.ExchangeFanout, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare exchange error: %+v", err)
	}
	err = ch.ExchangeDeclare("my_message_exchange", amqp091.ExchangeDirect, false, false, false, false, amqp091.Table{
		"alternate-exchange": "my_alternative_exchange",
	})
	if err != nil {
		log.Fatalf("declare exchange error: %+v", err)
	}

	defer func() {
		err = ch.ExchangeDelete("my_alternative_exchange", false, false)
		if err != nil {
			log.Printf("delete exchange error: %+v", err)
		}
		err = ch.ExchangeDelete("my_message_exchange", false, false)
		if err != nil {
			log.Printf("delete exchange error: %+v", err)
		}
	}()
	var routingKey string
	for i := 0; i < 40; i++ {
		switch i % 4 {
		case 0:
			routingKey = "route_a"
		case 1:
			routingKey = "route_c"
		case 2:
			routingKey = "route_e"
		case 3:
			routingKey = "route_g"

		}
		err = ch.Publish("my_message_exchange", routingKey, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("[%s]hello world %d!", routingKey, i)),
		})
		if err != nil {
			log.Fatalf("publish error: %+v", err)
		}
	}
	select {}
}

func alternativeConsumer(routingKey string) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}
	queue, err := ch.QueueDeclare(fmt.Sprintf("%s_queue", routingKey), false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}
	defer func() {
		_, err := ch.QueueDelete(queue.Name, false, false, false)
		if err != nil {
			log.Printf("delete queue error: %+v", err)
		}
	}()
	err = ch.QueueBind(queue.Name, routingKey, "my_message_exchange", false, nil)
	if err != nil {
		log.Fatalf("bind queue error: %+v", err)
	}

	defer func() {
		err := ch.QueueUnbind(queue.Name, routingKey, "my_message_exchange", nil)
		if err != nil {
			log.Printf("unbind queue error: %+v", err)
		}
	}()
	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume message error: %+v", err)
	}
	for msg := range messages {
		log.Printf("consume message: routingKey:%s,body:%s", msg.RoutingKey, string(msg.Body))
	}
}

func alternativeExchangeConsumer(routingKey string) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v", err)
	}
	queue, err := ch.QueueDeclare("my_alternative_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}
	defer func() {
		_, err := ch.QueueDelete(queue.Name, false, false, false)
		if err != nil {
			log.Printf("delete queue error: %+v", err)
		}
	}()
	err = ch.QueueBind(queue.Name, routingKey, "my_alternative_exchange", false, nil)
	if err != nil {
		log.Fatalf("bind queue error: %+v", err)
	}
	defer func() {
		err := ch.QueueUnbind(queue.Name, routingKey, "my_alternative_exchange", nil)
		if err != nil {
			log.Printf("unbind queue error: %+v", err)
		}
	}()

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume message error: %+v", err)
	}
	for msg := range messages {
		log.Printf("consume no route message: routingKey:%s,body:%s", msg.RoutingKey, string(msg.Body))
	}
}
