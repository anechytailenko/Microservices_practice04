package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewSubscriber(url, exchangeName, queueName, routingKey string) (*Subscriber, <-chan amqp.Delivery, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open channel: %w", err)
	}

	success := false
	defer func() {
		if !success {
			ch.Close()
			conn.Close()
		}
	}()

	// set up deadletter exchange and queue
	dlxName := "events.dlx"
	dlqName := "notifications.dlq"

	err = ch.ExchangeDeclare(dlxName, "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare DLX: %w", err)
	}

	_, err = ch.QueueDeclare(dlqName, true, false, false, false, amqp.Table{"x-queue-type": "quorum"})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare DLQ: %w", err)
	}

	err = ch.QueueBind(dlqName, "", dlxName, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to bind DLQ: %w", err)
	}

	// set up the topic  exchange (the same as in publisher) and bind the queue to it with message poisoning policy (5 retry -> deadletter exchange)
	err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare main exchange: %w", err)
	}

	args := amqp.Table{
		"x-queue-type":           "quorum",
		"x-dead-letter-exchange": dlxName,
		"x-delivery-limit":       int32(5),
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare main queue: %w", err)
	}

	err = ch.QueueBind(q.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to bind main queue: %w", err)
	}
	// fair dispatch
	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"notification_service",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Printf("[RabbitMQ] Subscribed to queue '%s' (binding: '%s')", queueName, routingKey)

	success = true

	return &Subscriber{conn: conn, ch: ch}, msgs, nil
}

func (s *Subscriber) Close() {
	if s.ch != nil {
		s.ch.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}
