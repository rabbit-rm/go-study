package rabbitmq

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestDirectExchange(t *testing.T) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673")
	if err != nil {
		log.Fatalf("create conn error:%+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	producerCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}
	const exchangeName = "logs_direct"
	err = producerCh.ExchangeDeclare(exchangeName, amqp091.ExchangeDirect, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare direct exchange error:%+v\n", err)
	}
	go func() {
		index := 1
		for {
			// waring info error
			var key string
			switch rand.Int31n(5) {
			case 0:
				key = "waring"
			case 4:
				key = "error"
			default:
				key = "info"
			}
			message := fmt.Sprintf("[%s] message_%d", key, index)
			if err := producerCh.Publish(exchangeName, key, false, false, amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			}); err != nil {
				log.Printf("send message error:%+v\n", err)
			}
			log.Printf("send message:%+v\n", message)
			index++
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// consumer
	consumerChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error:%+v\n", err)
	}
	err = consumerChannel.ExchangeDeclare(exchangeName, amqp091.ExchangeDirect, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare direct exchange error:%+v\n", err)
	}
	go func() {
		waringQueue, err := consumerChannel.QueueDeclare("", false, false, true, false, nil)
		if err != nil {
			log.Fatalf("declare waring queue error:%+v\n", err)
		}
		if err := consumerChannel.QueueBind(waringQueue.Name, "waring", exchangeName, false, nil); err != nil {
			log.Fatalf("bind queue error:%+v\n", err)
		}
		if err := consumerChannel.QueueBind(waringQueue.Name, "info", exchangeName, false, nil); err != nil {
			log.Fatalf("bind queue error:%+v\n", err)
		}
		consume, err := consumerChannel.Consume(waringQueue.Name, "", true, true, false, false, nil)
		if err != nil {
			log.Fatalf("create consumer error:%+v\n", err)
		}
		out, err := os.OpenFile("W:\\GoProject\\private\\study\\rabbitmq\\logs_waring.out", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("create out error:%+v\n", err)
		}
		for msg := range consume {
			_, _ = fmt.Fprintf(out, "receive message:%s\n", msg.Body)
		}
	}()
	go func() {
		infoQueue, err := consumerChannel.QueueDeclare("", false, false, true, false, nil)
		if err != nil {
			log.Fatalf("declare waring queue error:%+v\n", err)
		}
		if err := consumerChannel.QueueBind(infoQueue.Name, "info", exchangeName, false, nil); err != nil {
			log.Fatalf("bind queue error:%+v\n", err)
		}
		consume, err := consumerChannel.Consume(infoQueue.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("create consumer error:%+v\n", err)
		}
		out, err := os.OpenFile("W:\\GoProject\\private\\study\\rabbitmq\\logs_info.out", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("create out error:%+v\n", err)
		}
		for msg := range consume {
			_, _ = fmt.Fprintf(out, "receive message:%s\n", msg.Body)
		}
	}()

	go func() {
		errorQueue, err := consumerChannel.QueueDeclare("", false, false, true, false, nil)
		if err != nil {
			log.Fatalf("declare waring queue error:%+v\n", err)
		}
		if err := consumerChannel.QueueBind(errorQueue.Name, "error", exchangeName, false, nil); err != nil {
			log.Fatalf("bind queue error:%+v\n", err)
		}
		consume, err := consumerChannel.Consume(errorQueue.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("create consumer error:%+v\n", err)
		}
		out, err := os.OpenFile("W:\\GoProject\\private\\study\\rabbitmq\\logs_error.out", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("create out error:%+v\n", err)
		}
		for msg := range consume {
			_, _ = fmt.Fprintf(out, "receive message:%s\n", msg.Body)
		}
	}()

	select {}
}
