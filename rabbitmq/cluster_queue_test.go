package rabbitmq

import (
	"fmt"
	"log"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

func TestClusterQueue(t *testing.T) {

	var addrs = []string{
		"amqp://root:123456@192.168.204.132:5671",
		"amqp://root:123456@192.168.204.132:5671",
		"amqp://root:123456@192.168.204.132:5671",
	}
	conn, err := amqp091.Dial(addrs[1])
	if err != nil {
		log.Fatalf("create connetion error: %+v\n", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v\n", err)
	}

	queue, err := ch.QueueDeclare("ttl_queue", true, false, false, false, amqp091.Table{
		"x-message-ttl": 60000,
	})
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}

	fmt.Printf("queue:%s\n", queue.Name)
}
