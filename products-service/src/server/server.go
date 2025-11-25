package server

import (
	"fmt"
	"log"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/Adityadangi14/ecomm_ai/pkg/WDB"
	"github.com/Adityadangi14/ecomm_ai/products-service/handlers"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/llm"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/mq"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/repository"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/routes"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/schema"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

type Server struct {
	db   *WDB.WDB
	amqp *amqp.Connection
	cfg  *config.Config
}

func NewProductServer(wdb *WDB.WDB, mq *amqp.Connection, cfg *config.Config) *Server {
	return &Server{db: wdb, amqp: mq, cfg: cfg}
}

func (s *Server) Run() error {

	app := fiber.New()

	aiClient := llm.NewAiClient()

	proPub, err := mq.NewProductsPublisher(s.amqp, s.cfg, aiClient)

	err = schema.CreateProductClass(s.db.DB)

	if err != nil {
		fmt.Println("Failed to create product class", err)
	}

	if err != nil {
		return err
	}

	err = proPub.SetupExchangeAndQueue(s.cfg.RabbitMQ.Exchange,
		s.cfg.RabbitMQ.Queue,
		s.cfg.RabbitMQ.RoutingKey,
		s.cfg.RabbitMQ.ConsumerTag)

	if err != nil {
		fmt.Println("Failed to setup exchange  and queue", err)
	}

	//defer proPub.CloseChan()

	prodRepo := repository.NewProductRepository(s.db)
	prodConu := mq.NewProductsConsumer(s.amqp, prodRepo, aiClient)

	go func() {
		err := prodConu.StartConsumer(
			s.cfg.RabbitMQ.WorkerPoolSize,
			s.cfg.RabbitMQ.Exchange,
			s.cfg.RabbitMQ.Queue,
			s.cfg.RabbitMQ.RoutingKey,
			s.cfg.RabbitMQ.ConsumerTag,
		)

		if err != nil {
			fmt.Println("failed to start product consumer: ", err)

		}
	}()

	apiHandler := handlers.NewHandler(proPub, prodRepo, aiClient)

	routes.RegisterRoutes(app, *apiHandler)

	log.Fatal(app.Listen(":3000"))

	return nil
}
