package mq

import (
	"context"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	logg   logger.Logg
	conf   *config.Config
	conn   *amqp.Connection
	chanel *amqp.Channel
}

type RabbitAPI interface {
	Connect(ctx context.Context) error
	Publish(exchange string, queue string, body []byte) error
	Consume(ctx context.Context, queue string, consumer string) (<-chan []byte, error)
	Close() error
}

func New(logger logger.Logg, config *config.Config) RabbitAPI {
	return &Rabbit{
		logg: logger,
		conf: config,
	}
}

func (r *Rabbit) Connect(ctx context.Context) error {
	var err error
	r.conn, err = amqp.Dial(r.conf.RabbitMQ.ConnectString)
	if err != nil {
		r.logg.Error("failed to work with RabbitMQ ", err)
		return err
	}
	r.chanel, err = r.conn.Channel()
	if err != nil {
		r.logg.Error("failed to open a channel ", err)
		r.conn.Close()
		return err
	}
	return nil
}

func queueDeclare(chanel *amqp.Channel, queue string) {
	chanel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *Rabbit) Publish(exchange string, queue string, body []byte) error {
	queueDeclare(r.chanel, queue)
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	err := r.chanel.PublishWithContext(ctx,
		exchange, // exchange
		queue,    // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		r.logg.Error("failed to publish the message ", err)
		return err
	}
	return nil
}

func (r *Rabbit) Consume(ctx context.Context, queue string, consumer string) (<-chan []byte, error) {
	queueDeclare(r.chanel, queue)
	msgs, err := r.chanel.Consume(
		queue,    // queue
		consumer, // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		r.logg.Error("failed to consume the message ", err)
		return nil, err
	}
	messages := make(chan []byte)
	go func() {
		defer close(messages)
		for {
			select {
			case <-ctx.Done():
				return
			case del, ok := <-msgs:
				if !ok {
					return
				}
				if err = del.Ack(false); err != nil {
					r.logg.Error("confirming the message error: ", err)
				}
				select {
				case <-ctx.Done():
					return
				case messages <- del.Body:
				}
			}
		}
	}()
	return messages, nil
}

func (r *Rabbit) Close() error {
	if err := r.chanel.Close(); err != nil {
		r.logg.Error("channel connection error: ", err)
		return err
	}
	if err := r.conn.Close(); err != nil {
		r.logg.Error("rabbitMQ connection error: ", err)
		return err
	}
	r.logg.Info("connect to rabbitMQ is closed")
	return nil
}
