package rabbitmq

import (
	"log"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

func TestMessageTransaction(t *testing.T) {
	messageTransactionPublisher()
}

func messageTransactionPublisher() {
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
	err = ch.Tx()
	if err != nil {
		log.Fatalf("set tx mode error: %+v", err)
	}

	queue, err := ch.QueueDeclare("test_tx_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v", err)
	}

	for i := 0; i < 10; i++ {
		err = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte("message"),
		})
		if err != nil {
			_ = ch.TxRollback()
			log.Fatalf("publish message error: %+v", err)
		}
		_ = ch.TxCommit()
	}

}
