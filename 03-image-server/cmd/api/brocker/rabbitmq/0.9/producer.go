package producer

import (
	"context"
	"errors"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrFailedPublish     = errors.New("Failed to publish message.")
	ErrUnexpectedOutcome = errors.New("Unexpected message outcome.")
)

type RabbitProducer struct {
	logger  *log.Logger
	ctx     context.Context
	queue   amqp.Queue
	channel *amqp.Channel
}

func NewRabbitProducer(brokerURI string) *RabbitProducer {
	producer := new(RabbitProducer)
	producer.logger = log.New(os.Stdout, "rabbitmq: ", log.Ldate|log.Ltime)
	ctx := context.Background()
	producer.ctx = ctx

	conn, err := amqp.Dial(brokerURI)
	if err != nil {
		producer.logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		producer.logger.Fatalf("Failed to establish connection: %v", err)
	}

	queueName := os.Getenv("QUEUE_NAME")
	producer.logger.Printf("Queue name is %s\n", queueName)
	q, err := ch.QueueDeclare(
		queueName,
		true,  // durability
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		producer.logger.Fatalf("Failed to declare queue: %v", err)
	}

	producer.channel = ch
	producer.queue = q

	producer.logger.Println("Successful creation of producer.")
	return producer
}

func (producer *RabbitProducer) Send(body string) error {
	err := producer.channel.PublishWithContext(producer.ctx,
		"",                  // exchange
		producer.queue.Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	return err
}
