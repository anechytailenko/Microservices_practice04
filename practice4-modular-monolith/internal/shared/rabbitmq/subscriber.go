package rabbitmq

import (
	"fmt"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewSubscriber(url, exchangeName, queueName string, bindingKeys []string, dlxName, dlqName, consumerName string) (*Subscriber, <-chan amqp.Delivery, error) {
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

	err = ch.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare DLX: %w", err)
	}

	_, err = ch.QueueDeclare(dlqName, true, false, false, false, amqp.Table{"x-queue-type": "quorum"})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// use dlq name as the routing to its queue
	err = ch.QueueBind(dlqName, dlqName, dlxName, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to bind DLQ: %w", err)
	}

	err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare main exchange: %w", err)
	}

	args := amqp.Table{
		"x-queue-type":              "quorum",
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": dlqName,
		"x-delivery-limit":          int32(5),
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare main queue: %w", err)
	}

	for _, key := range bindingKeys {
		err = ch.QueueBind(q.Name, key, exchangeName, false, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to bind main queue to key %s: %w", key, err)
		}
	}

	// Fair dispatch
	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start consuming: %w", err)
	}

	keysLog := strings.Join(bindingKeys, ", ")
	log.Printf("[RabbitMQ] Subscribed to queue '%s' (bindings: [%s])", queueName, keysLog)

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
