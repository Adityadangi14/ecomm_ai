package mq

import (
	"time"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/llm"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type ProductPublisher interface {
	SetupExchangeAndQueue(exchange, queueName, bindingKey, consumerTag string) error
	Publish(body []byte, contentType string) error
	CloseChan() error
}

var _ ProductPublisher = (*Productpublisher)(nil)

type Productpublisher struct {
	amqpChan *amqp.Channel
	cfg      *config.Config
	Aiclient llm.Aiclient
}

func NewProductsPublisher(mqConn *amqp.Connection, cfg *config.Config, aiClient llm.Aiclient) (*Productpublisher, error) {

	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}
	return &Productpublisher{amqpChan: amqpChan, cfg: cfg, Aiclient: aiClient}, nil
}

func (p *Productpublisher) SetupExchangeAndQueue(exchange, queueName, bindingKey, consumerTag string) error {
	err := p.amqpChan.ExchangeDeclare(
		exchange,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)

	if err != nil {
		return err
	}

	queue, err := p.amqpChan.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)

	err = p.amqpChan.QueueBind(
		queue.Name,
		bindingKey,
		exchange,
		queueNoWait,
		nil,
	)

	if err != nil {
		return errors.Wrap(err, "Error ch.QueueBind")
	}
	return nil

}

func (p *Productpublisher) CloseChan() error {
	if err := p.amqpChan.Close(); err != nil {
		return err
	}
	return nil
}

func (p *Productpublisher) Publish(body []byte, contentType string) error {
	if err := p.amqpChan.Publish(
		p.cfg.RabbitMQ.Exchange,
		p.cfg.RabbitMQ.RoutingKey,
		publishMandatory,
		publishImmediate,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		return errors.Wrap(err, "ch.Publish")
	}
	return nil
}
