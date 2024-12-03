package event

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type Payload struct {
	Topic string         `json:"topic"`
	Data  map[string]any `json:"data"`
}

func NewConsumer(conn *amqp.Connection, queueName string) (*Consumer, error) {
	consumer := &Consumer{
		conn:      conn,
		queueName: queueName,
	}

	err := consumer.setup()
	if err != nil {
		return &Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel, consumer.queueName)
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	queue, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if err := ch.QueueBind(
			queue.Name,         // queue name
			topic,              // routing key
			consumer.queueName, // exchange
			false,              // no-wait
			nil,                // arguments
		); err != nil {
			return err
		}
	}

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)

		}
	}()

	fmt.Printf("Consumer is ready, waiting for message [Exchange, Queue] [%s, %s]\n", consumer.queueName, queue.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Topic {
	case "notify":
		err := notify(payload)
		if err != nil {
			fmt.Println("Failed to notify")
		}
	}
}

func notify(payload Payload) error {
	return nil
}
