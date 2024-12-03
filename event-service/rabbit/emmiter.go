package rabbit

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emmiter struct {
	connection *amqp.Connection
}

func (e *Emmiter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel, "notify")
}

func (e *Emmiter) Push(event, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	log.Println("Pushing to channel")

	err = channel.Publish(
		"notify", // exchange
		severity, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		})
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (*Emmiter, error) {
	emitter := &Emmiter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return &Emmiter{}, err
	}

	return emitter, nil
}
