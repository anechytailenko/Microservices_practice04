package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	confirms chan amqp.Confirmation
	exchange string
}

func NewPublisher(url string, exchangeName string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set confirm mode: %w", err)
	}

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	logger.Printf(context.Background(), "[RabbitMQ] Connected and configured exchange '%s'", exchangeName)

	return &Publisher{
		conn:     conn,
		ch:       ch,
		confirms: confirms,
		exchange: exchangeName,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, routingKey string, payload []byte) error {

	corrID := ctxutil.GetCorrelationID(ctx)
	err := p.ch.PublishWithContext(ctx,
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Headers: amqp.Table{
				"X-Correlation-Id": corrID,
			},
			Body: payload,
		})
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	select {
	case confirmed := <-p.confirms:
		if !confirmed.Ack {
			return fmt.Errorf("broker NACKed the message")
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for publisher confirmation")
	}
}

func (p *Publisher) Close() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
