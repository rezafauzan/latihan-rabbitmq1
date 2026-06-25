package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct{
	RMQ *amqp091.Connection
}

func NewRabbit() (RabbitMQ, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
		return RabbitMQ{}, err
	}
	defer conn.Close()

	return RabbitMQ{
		conn,
	}, nil
}

func ProduceMessage(connection *amqp091.Connection, queueName string, messages string) error {
	ch, err := connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(fmt.Sprintf("Queue %s", queueName), true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp091.Publishing{ContentType: "text/plain", Body: []byte(messages)})
	if err != nil {
		return err
	}

	return nil
}