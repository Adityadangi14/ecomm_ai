package main

import (
	"log"
	"os"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/Adityadangi14/ecomm_ai/log-service/src/server"
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

	s := server.NewServer(cfg)

	err = s.Run()

	if err != nil {
		log.Fatalf("failed to run log server %v", err)
	}
}
