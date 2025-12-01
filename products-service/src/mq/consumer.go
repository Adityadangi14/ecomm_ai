package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/llm"
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

	prefetchCount  = 24
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
	Aiclient llm.Aiclient
}

func NewProductsConsumer(ampqConn *amqp.Connection, prodRep repository.ProductRepository, aiClient llm.Aiclient) *ProductConsumer {
	return &ProductConsumer{amqpConn: ampqConn, prodRepo: prodRep, Aiclient: aiClient}
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

func (p *ProductConsumer) worker(ctx context.Context, id int, jobs <-chan amqp.Delivery) {
	for delivery := range jobs {
		// fmt.Printf("Worker %d processing: %s\n", id, delivery.Body)

		var body models.Product
		if err := json.Unmarshal(delivery.Body, &body); err != nil {
			fmt.Printf("Worker %d: invalid JSON: %v\n", id, err)
			_ = delivery.Reject(false)
			continue
		}

		res, err := p.Aiclient.ProcessProduct(body)

		if err != nil {

			fmt.Println("error in processing product", err)

			_ = delivery.Reject(true)
		} else {
			fmt.Println("product to save ", res)
			err = p.prodRepo.SaveProduct(ctx, res)
			if err != nil {
				fmt.Printf("Worker %d: save failed: %v\n", id, err)
				_ = delivery.Reject(true)
				continue
			}

			if err := delivery.Ack(false); err != nil {
				fmt.Printf("Worker %d: Ack failed: %v\n", id, err)
			}
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

	jobs := make(chan amqp.Delivery, workerPoolSize*2)

	// Start worker pool
	for i := 0; i < workerPoolSize; i++ {
		go p.worker(ctx, i, jobs)
	}

	// Consumer loop
	go func() {
		for d := range deliveries {
			jobs <- d
		}
		close(jobs)
	}()
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))

	return chanErr
}
