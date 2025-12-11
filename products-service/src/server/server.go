package server

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/Adityadangi14/ecomm_ai/pkg/WDB"
	"github.com/Adityadangi14/ecomm_ai/pkg/redis"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/handlers"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/llm"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/logging"
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

	url := os.Getenv("logging_url")

	serviceName := "products-service"

	logging := logging.NewLogging(url, serviceName)
	logHandler := slog.New(logging).With(slog.String("source", "enabled"))
	slog.SetDefault(logHandler)

	app := fiber.New()

	rdb, err := redis.ConnectToRedis(s.cfg)

	if err != nil {
		return fmt.Errorf("failed to connect to redis:%v", err)
	}

	prodRepo := repository.NewProductRepository(s.db)

	aiClient := llm.NewAiClient(rdb, prodRepo)

	proPub, err := mq.NewProductsPublisher(s.amqp, s.cfg, aiClient)

	if err != nil {
		return err
	}

	err = schema.CreateProductClass(s.db.DB)

	if err != nil {
		slog.Error("Failed to create product class", "error", err)
	}

	err = proPub.SetupExchangeAndQueue(s.cfg.RabbitMQ.Exchange,
		s.cfg.RabbitMQ.Queue,
		s.cfg.RabbitMQ.RoutingKey,
		s.cfg.RabbitMQ.ConsumerTag)

	if err != nil {
		slog.Error("Failed to setup exchange  and queue", "error", err)
	}

	//defer proPub.CloseChan()

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
			slog.Error("failed to start product consumer: ", "error", err)

		}
	}()

	apiHandler := handlers.NewHandler(proPub, prodRepo, aiClient, rdb)

	routes.RegisterRoutes(app, *apiHandler)
	slog.Info("Product service server started at :3000")
	log.Fatal(app.Listen(":3000"))

	return nil
}
