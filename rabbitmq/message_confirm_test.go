package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestPublisherConfirm(t *testing.T) {
	// messageConfirmPublisher()
	// messageConfirmPublisherMultiple()
	messageConfirmPublisherAsync()
}

func messageConfirmPublisher() {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
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

	err = ch.Confirm(false)
	if err != nil {
		log.Fatalf("set confirm mode error: %+v", err)
	}

	queue, err := ch.QueueDeclare("test_confirm_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}
	start := time.Now()
	for i := 0; i < 10; i++ {
		// 500微妙
		/*ch.Publish("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("message %d", i)),
		})*/
		// 7ms
		confirm, err := ch.PublishWithDeferredConfirm("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("message %d", i)),
		})
		if err != nil {
			log.Fatalf("publish message error: %+v", err)
		}
		if confirm.Wait() {
			log.Printf("message %d published,id:%d", i, confirm.DeliveryTag)
		}
	}
	log.Printf("publish 10 message cost: %v", time.Since(start))
}

func messageConfirmPublisherMultiple() {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
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

	err = ch.Confirm(false)
	if err != nil {
		log.Fatalf("set confirm mode error: %+v", err)
	}

	queue, err := ch.QueueDeclare("test_confirm_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}
	start := time.Now()
	batchSize := 5
	noConfirmMesNum := 0
	for i := 0; i < 10; i++ {
		confirm, err := ch.PublishWithDeferredConfirm("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("message %d", i)),
		})
		if err != nil {
			log.Fatalf("publish message error: %+v", err)
		}
		noConfirmMesNum++
		if noConfirmMesNum == batchSize {
			if confirm.Wait() {
				log.Printf("message %d published,id:%d", i, confirm.DeliveryTag)
			}
		}
	}
	log.Printf("publish 10 message cost: %v", time.Since(start))
}

func messageConfirmPublisherAsync() {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
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

	err = ch.Confirm(false)
	if err != nil {
		log.Fatalf("set confirm mode error: %+v", err)
	}

	queue, err := ch.QueueDeclare("test_confirm_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}

	ch1 := make(chan amqp091.Confirmation)
	ch2 := make(chan amqp091.Return)

	ch.NotifyPublish(ch1)
	ch.NotifyReturn(ch2)

	go func() {
		for c := range ch1 {
			log.Printf("message published,id:%d,isAck:%v\n", c.DeliveryTag, c.Ack)
		}
	}()
	go func() {
		for c := range ch2 {
			log.Printf("message publish failed,body:%s\n", c.Body)
		}
	}()
	start := time.Now()
	for i := 0; i < 10; i++ {
		_, _ = ch.PublishWithDeferredConfirm("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("message %d", i)),
		})
	}
	log.Printf("publish 10 message cost: %v", time.Since(start))
}
