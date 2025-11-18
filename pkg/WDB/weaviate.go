package WDB

import (
	"context"
	"fmt"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type WDB struct {
	DB *weaviate.Client
}

func NewWeaviateDB(c *config.Config) (*WDB, error) {

	cfg := weaviate.Config{
		Scheme: c.Weaviate.Scheme,
		Host:   c.Weaviate.Host,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("Weaviate client is nil")
	}
	_, err = client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ReadyCheck Failed %v", err)
	}
	return &WDB{DB: client}, nil

}
