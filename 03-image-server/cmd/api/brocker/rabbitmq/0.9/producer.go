package producer

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrFailedPublish     = errors.New("Failed to publish message.")
	ErrUnexpectedOutcome = errors.New("Unexpected message outcome.")
)

type MessageHandler interface {
	Handle(body []byte) error
}

type RabbitProducer struct {
	logger     *log.Logger
	ctx        context.Context
	queue      amqp.Queue
	channel    *amqp.Channel
	msgHandler MessageHandler
	pending    map[string]struct{}
	mu         sync.Mutex
}

func NewRabbitProducer(brokerURI string, msgHandler MessageHandler) *RabbitProducer {
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
	producer.logger.Printf("Queue name --- %s\n", queueName)
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
	producer.msgHandler = msgHandler
	producer.pending = make(map[string]struct{})

	// handling replies
	go producer.onMessage()

	producer.logger.Println("Successful creation of producer.")
	return producer
}

func (producer *RabbitProducer) Send(body string) error {
	correlationID := uuid.New().String()

	producer.mu.Lock()
	producer.pending[correlationID] = struct{}{} // ← save it
	producer.mu.Unlock()

	return producer.channel.PublishWithContext(producer.ctx,
		"",
		producer.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			ReplyTo:       os.Getenv("REPLY_QUEUE_NAME"),
			CorrelationId: correlationID, // ← use saved ID
			Body:          []byte(body),
		})
}

func (producer *RabbitProducer) onMessage() error {
	queueName := os.Getenv("REPLY_QUEUE_NAME")
	replyQueue, err := producer.channel.QueueDeclare(
		queueName, false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	msgs, err := producer.channel.Consume(
		replyQueue.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	producer.logger.Println("Starting to receive answers...")
	for msg := range msgs {
		producer.mu.Lock()
		_, ok := producer.pending[msg.CorrelationId]
		if ok {
			delete(producer.pending, msg.CorrelationId) // ← clean up
		}
		producer.mu.Unlock()

		if ok {
			producer.logger.Printf("Received answer: %s", msg.Body)
			producer.msgHandler.Handle(msg.Body)
		}
	}

	return nil
}
