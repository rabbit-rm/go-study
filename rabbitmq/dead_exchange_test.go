package rabbitmq

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestDLX(t *testing.T) {
	conn, err := amqp091.Dial(rabbitUrl)
	if err != nil {
		log.Fatalf("create connection error: %+v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	// declare dlx
	var dlxName = "message-dlx"
	{
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}

		if err = ch.ExchangeDeclare(
			dlxName,
			"direct",
			true,
			false,
			false,
			false,
			nil); err != nil {
			log.Fatalf("declare exchange error: %+v", err)
		}

		queue, err := ch.QueueDeclare(
			"dlx-queue",
			false,
			true,
			true,
			false,
			nil)
		if err != nil {
			log.Fatalf("declare queue error: %+v", err)
		}

		if err := ch.QueueBind(queue.Name, "dlx-queue-key-1", dlxName, false, nil); err != nil {
			log.Fatalf("bind queue error: %+v", err)
		}

		messages, err := ch.Consume(
			queue.Name,
			"",
			true,
			true,
			false,
			false,
			nil)
		if err != nil {
			log.Fatalf("consume error: %+v", err)
		}

		go func() {
			for msg := range messages {
				var buf bytes.Buffer
				buf.WriteString("dlx message:\n")
				buf.WriteString("\theaders:\n")
				headers := msg.Headers
				for k, v := range headers {
					buf.WriteString(fmt.Sprintf("\t\t%s:%v\n", k, v))
				}
				buf.WriteString(fmt.Sprintf("\tbody:%s\n", msg.Body))
				log.Printf("dlx message: %s", buf.String())
			}
		}()
	}

	{
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("create channel error: %+v", err)
		}

		queue, err := ch.QueueDeclare(
			"ttl-queue",
			false,
			true,
			true,
			false,
			amqp091.Table{
				"x-dead-letter-routing-key": "dlx-queue-key-1",
				"x-dead-letter-exchange":    dlxName,
				"x-max-length":              5,
			})
		if err != nil {
			log.Fatalf("declare queue error: %+v", err)
		}

		go func() {
			index := 1
			for {
				_ = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("message-%d", index)),
					Expiration:  strconv.Itoa(int(rand.Int31n(10) * 100)),
				})
				time.Sleep(100 * time.Millisecond)
				index++
			}
		}()
	}

	select {}

}
