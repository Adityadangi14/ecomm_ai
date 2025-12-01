package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/llm"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/mq"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/repository"
	"github.com/Adityadangi14/ecomm_ai/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	ProductHandlers ProductHandlers
	QueryHandler    QueryHandler
}

func NewHandler(pub mq.ProductPublisher, productRepo repository.ProductRepository, aiCLient llm.Aiclient, rdb *redis.Client) *Handlers {
	prodHandler := NewProductHandlers(pub, productRepo)
	queryHandler := NewQueryHandler(aiCLient, rdb)
	return &Handlers{
		ProductHandlers: prodHandler,
		QueryHandler:    queryHandler,
	}

}

type ProductHandlers interface {
	UploadProducts(c *fiber.Ctx) error
	DeleteAllProducts(c *fiber.Ctx) error
}

type prodHandlers struct {
	prodPublisher mq.ProductPublisher
	productRepo   repository.ProductRepository
}

func NewProductHandlers(pub mq.ProductPublisher, productRepo repository.ProductRepository) ProductHandlers {
	return &prodHandlers{prodPublisher: pub, productRepo: productRepo}
}

func (p *prodHandlers) UploadProducts(c *fiber.Ctx) error {

	var products models.ProductsModel

	if err := c.BodyParser(&products); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, fmt.Sprintf("failed to parse body: %v", err))
	}

	var errors []struct {
		ProductID string
		err       error
	}

	for _, prod := range products.Products {

		byt, err := json.Marshal(prod)

		if err != nil {
			return utils.Fail(c, fiber.StatusBadRequest, fmt.Sprintf("failed to marshal body: %v", err))
		}
		err = p.prodPublisher.Publish(byt, "text")

		if err != nil {
			fmt.Println(err)
			errors = append(errors, struct {
				ProductID string
				err       error
			}{
				ProductID: prod.ID,
				err:       err,
			})
		}

	}

	if len(errors) != 0 {
		return utils.Fail(c, fiber.StatusMultiStatus, fmt.Sprintf("upload failed for following items %v", errors))
	}

	utils.Success(c, "Product upload has been started successfully")

	return nil
}

func (p *prodHandlers) DeleteAllProducts(c *fiber.Ctx) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	err := p.productRepo.DeleteAllProducts(ctx)

	if err != nil {
		return utils.Fail(c, fiber.StatusInternalServerError, fmt.Sprintf("unable to delete products.%v", err))
	}
	return utils.Success(c, "Successfully deleted products")

}
