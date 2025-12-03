package handlers

import "github.com/gofiber/fiber/v2"

type Handler struct {
	LogHandler LogHandler
}

func NewHandler(logHandler LogHandler) Handler {
	return Handler{LogHandler: logHandler}
}

type LogHandler interface {
	HandleLogs(*fiber.Ctx) error
}

type logHandler struct {
}

func NewLogHandler() LogHandler {
	return &logHandler{}
}

func (l *logHandler) HandleLogs(c *fiber.Ctx) error {
	return nil
}
