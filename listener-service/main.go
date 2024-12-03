package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	consumer, err := event.NewConsumer(rabbitConn, "notify")
	if err != nil {
		panic(err)
	}

	topics := []string{"event_add", "event_update", "event_remove", "event_register", "user_registered"}
	err = consumer.Listen(topics)
	if err != nil {
		panic(err)
	}
}
func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ Not yet ready...")
			counts++
		} else {
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Retrying in %v\n", backOff)
		time.Sleep(backOff)
	}

	return connection, nil
}
