package main

import (
	"log"
	"os"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/Adityadangi14/ecomm_ai/pkg/WDB"
	"github.com/Adityadangi14/ecomm_ai/pkg/rabbitmq"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/server"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.GetConfig(os.Getenv("config"))
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConn.Close()

	wdb, err := WDB.NewWeaviateDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewProductServer(wdb, amqpConn, cfg)

	err = s.Run()

	if err != nil {
		log.Fatal(err)
	}
}
