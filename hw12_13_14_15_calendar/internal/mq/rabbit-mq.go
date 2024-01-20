package mq

import (
	"context"
	"sync"
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
	Connect(ctx context.Context, wg *sync.WaitGroup) error
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

func (r *Rabbit) Connect(ctx context.Context, wg *sync.WaitGroup) error {
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
	wg.Add(1)
	go func() {
		<-ctx.Done()
		if err = r.Close(); err != nil {
			r.logg.Error(err)
		} else {
			r.logg.Info("connect to rabbitMQ is closed")
		}
		wg.Done()
	}()
	return nil
}

func (r *Rabbit) Publish(exchange string, queue string, body []byte) error {
	var once sync.Once
	f := func() {
		r.chanel.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
	}
	once.Do(f)

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
	var once sync.Once
	f := func() {
		r.chanel.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
	}
	once.Do(f)

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
		defer func() {
			close(messages)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case del, ok := <-msgs:
				if !ok {
					return
				}
				if err = del.Ack(false); err != nil {
					r.logg.Error(err)
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
	err := r.chanel.Close()
	if err != nil {
		return err
	}
	err = r.conn.Close()
	return err
}
