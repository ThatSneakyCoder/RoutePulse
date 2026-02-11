package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const AnalyticsExchange = "analytics.events" // exchange name should be some purpose based

type MessageHandler func(context.Context, amqp.Delivery) error

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    []byte `json:"data"`
}

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}

	rmq := &RabbitMQ{
		conn:    conn,
		Channel: ch,
	}

	if err := rmq.setupExchangesAndQueues(); err != nil {
		// clean up if setup fails
		rmq.Close()
		return nil, fmt.Errorf("failed to setup exchanges and queues: %v", err)
	}

	return rmq, nil
}

func (r *RabbitMQ) setupExchangesAndQueues() error {
	// 1. create an exchange
	err := r.Channel.ExchangeDeclare(
		AnalyticsExchange, // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)

	if err != nil {
		log.Fatalf("failed to create exchange: %v", err)
		return fmt.Errorf("failed to declare exchange %s: %v", AnalyticsExchange, err)
	}

	if err := r.declareAndBindQueue(
		AnalyticsIdentityQueue,
		[]string{
			IdentityUserRegisteredEvent,
			IdentityUserEmailVerifiedEvent,
			OrganizationCreatedEvent,
			OrganizationMemberAddedEvent,
			IdentityUserLoggedInEvent,
		},
		AnalyticsExchange,
	); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) declareAndBindQueue(queueName string, messageTypes []string, exchange string) error {
	q, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no wait
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, msg := range messageTypes {
		if err := r.Channel.QueueBind(
			q.Name,   // queue name
			msg,      // routing key
			exchange, // exchange
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue to %s: %v", queueName, err)
		}
	}

	return nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}

	if r.Channel != nil {
		r.Channel.Close()
	}
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message AmqpMessage) error {
	log.Printf("Publishing message with routing key: %s", routingKey)

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	return r.Channel.PublishWithContext(ctx,
		AnalyticsExchange, // exchange
		routingKey,        // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         jsonMsg,
			DeliveryMode: amqp.Persistent,
		})
}

func (r *RabbitMQ) ConsumeMessages(queueName string, handler MessageHandler) error {

	err := r.Channel.Qos(
		1,     // limit to 1 acknowledged message per consumer
		0,     // no specific limit on message size
		false, // app prefetchcount to each consumer individually
	)

	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	// start delivering the messages from the specified queue to me
	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatal("failed to consume message: %w", err)
		return err
	}

	ctx := context.Background()

	go func() {
		for msg := range msgs {
			log.Printf(" Received a message %s", msg.Body)

			if err := handler(ctx, msg); err != nil {
				log.Fatalf("failed to handle the message: %v", err)
			}
			// Only ack if the handler succeeds
			_ = msg.Ack(false)
		}
	}()

	return nil
}
