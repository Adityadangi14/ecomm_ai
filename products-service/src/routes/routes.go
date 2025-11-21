package routes

import (
	"github.com/Adityadangi14/ecomm_ai/products-service/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handlers handlers.Handlers) {
	v1 := app.Group("api/v1")

	v1.Post("/uploadProducts", handlers.ProductHandlers.UploadProducts)
	v1.Delete("/deleteAllProducts", handlers.ProductHandlers.DeleteAllProducts)
}
