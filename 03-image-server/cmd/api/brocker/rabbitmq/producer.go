package producer

import (
	"context"
	"errors"
	"log"
	"os"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

var (
	ErrFailedPublish     = errors.New("Failed to publish message.")
	ErrUnexpectedOutcome = errors.New("Unexpected message outcome.")
)

type RabbitProducer struct {
	logger    *log.Logger
	ctx       context.Context
	publisher *rmq.Publisher
}

func NewRabbitProducer(brokerURI string) *RabbitProducer {
	producer := new(RabbitProducer)
	producer.logger = log.New(os.Stdout, "rabbitmq: ", log.Ldate|log.Ltime)
	ctx := context.Background()
	producer.ctx = ctx

	env := rmq.NewEnvironment(brokerURI, nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		producer.logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	queueName := os.Getenv("QUEUE_NAME")
	_, err = conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: queueName})
	if err != nil {
		producer.logger.Fatalf("Failed to declare a queue: %v", err)
	}

	publisher, err := conn.NewPublisher(ctx, &rmq.QueueAddress{Queue: queueName}, nil)
	if err != nil {
		producer.logger.Fatalf("Failed to create a publisher: %v", err)
	}

	producer.publisher = publisher

	producer.logger.Println("Successful creation of producer.")
	return producer
}

func (producer *RabbitProducer) Send(body string) error {
	res, err := producer.publisher.Publish(producer.ctx, rmq.NewMessage([]byte(body)))
	if err != nil {
		return ErrFailedPublish
	}

	switch res.Outcome.(type) {
	case *rmq.StateAccepted:
	default:
		return ErrUnexpectedOutcome
	}

	producer.logger.Printf("[x] Sent\n%s\n", body)
	return nil
}
