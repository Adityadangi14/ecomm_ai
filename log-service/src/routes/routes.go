package routes

import (
	"github.com/Adityadangi14/ecomm_ai/log-service/src/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler handlers.Handler) {
	g := app.Group("api/v1/")

	g.Post("/log", handler.LogHandler.HandleLogs)
}
