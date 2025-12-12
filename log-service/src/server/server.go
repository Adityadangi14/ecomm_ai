package server

import (
	"log"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/Adityadangi14/ecomm_ai/log-service/src/handlers"
	"github.com/Adityadangi14/ecomm_ai/log-service/src/kafka"
	"github.com/Adityadangi14/ecomm_ai/log-service/src/routes"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	cfg *config.Config
}

func NewServer(c *config.Config) *Server {
	return &Server{cfg: c}
}

func (s *Server) Run() error {

	logProducer := kafka.NewLogProducer(s.cfg)

	logHandler := handlers.NewLogHandler(logProducer)

	handler := handlers.NewHandler(logHandler)

	app := fiber.New()

	routes.RegisterRoutes(app, handler)

	log.Fatal(app.Listen(":4000"))

	return nil
}
