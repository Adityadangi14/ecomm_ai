package rabbitmq

import (
	"fmt"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/streadway/amqp"
)

// Initialize new RabbitMQ connection
func NewRabbitMQConn(cfg *config.Config) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	fmt.Println("conn adder", connAddr)
	return amqp.Dial(connAddr)
}
