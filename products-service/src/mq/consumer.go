package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/repository"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	publishMandatory = false
	publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type ProductConsumer struct {
	amqpConn *amqp.Connection
	prodRepo repository.ProductRepository
}

func NewProductsConsumer(ampqConn *amqp.Connection, prodRep repository.ProductRepository) *ProductConsumer {
	return &ProductConsumer{amqpConn: ampqConn, prodRepo: prodRep}
}

func (p *ProductConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueDeclare")
	}

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error  ch.Qos")
	}

	return ch, nil
}

func (p *ProductConsumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {
	var body models.ProductsModel
	for delivery := range messages {
		err := json.Unmarshal(delivery.Body, &body)
		if err != nil {
			delivery.Reject(false)
		} else {
			err = p.prodRepo.SaveProduct(ctx, body)

			if err != nil {
				delivery.Reject(true)
			}
			err := delivery.Ack(false)

			fmt.Println("failed to ack delivery", err)
		}

	}
}

func (p *ProductConsumer) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := p.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	for i := 0; i < workerPoolSize; i++ {
		go p.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))

	return chanErr
}
