package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestTopicExchange(t *testing.T) {
	// go topicPublisher("", "empty_message")
	// go topicPublisher("kern.info", "kern_info_message")
	// go topicPublisher("kern.waring", "kern_waring_message")
	// go topicPublisher("syslog.error", "syslog_error_message")
	// go topicPublisher("gin.error", "gin_error_message")
	// // 匹配任意消息
	// go topicConsumer("#")
	// // 没有匹配的 binging key，消息被丢弃
	// go topicConsumer("*")
	// // 匹配所有的 error 消息
	// go topicConsumer("*.error")
	// // 匹配所有的 kern 消息
	// go topicConsumer("kern.*")

	// go topicPublisher("", "empty_message")
	// go topicPublisher("..", "point_point_message")
	// // 没有匹配的 binging key，消息被丢弃
	// go topicConsumer("*")
	// // 可以匹配 routing key 为 ".." 消息
	// go topicConsumer("#.*")

	go topicPublisher("a", "a_message")
	go topicPublisher("a.", "a._message")
	// go topicPublisher("a.b", "a.b_message")
	// go topicPublisher("a.b.", "a.b._message")
	// go topicPublisher("a.b.c", "a.b.c_message")

	go topicConsumer("a.*.#")
	// go topicConsumer("a.#")

	select {}
}

// <facility>.<severity>
// <设施>.<严重性>
func topicPublisher(routingKey, msgPrefix string) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("create conn error: %+v\n", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v\n", err)
	}
	err = ch.ExchangeDeclare("logs_topic", amqp091.ExchangeTopic, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare exchange error: %+v\n", err)
	}
	index := 1
	for {
		message := fmt.Sprintf("%s_%d", msgPrefix, index)
		if err := ch.Publish("logs_topic", routingKey, false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		}); err != nil {
			log.Printf("publish message error: %+v\n", err)
			continue
		}
		index++
		time.Sleep(200 * time.Millisecond)
	}
}

func topicConsumer(routingKey string) {
	conn, err := amqp091.Dial("amqp://guest:guest@192.168.204.131:5673/")
	if err != nil {
		log.Fatalf("create conn error: %+v\n", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel error: %+v\n", err)
	}
	queue, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("declare queue error: %+v\n", err)
	}
	err = ch.QueueBind(queue.Name, routingKey, "logs_topic", false, nil)
	if err != nil {
		log.Fatalf("bind queue error: %+v\n", err)
	}
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	suffix := routingKey
	if strings.Contains(suffix, "*") {
		suffix = strings.ReplaceAll(suffix, "*", "all")
	}
	out, err := os.OpenFile(fmt.Sprintf("topic_%s.log", suffix), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("open file error: %+v\n", err)
	}
	for msg := range msgs {
		_, _ = fmt.Fprintf(out, "receive message,routingKey: %s, message: %s\n", routingKey, msg.Body)
	}
}
