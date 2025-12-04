package handlers

import (
	"encoding/json"

	"github.com/Adityadangi14/ecomm_ai/log-service/src/kafka"
	"github.com/Adityadangi14/ecomm_ai/log-service/src/models"
	"github.com/Adityadangi14/ecomm_ai/utils"
	"github.com/gofiber/fiber/v2"
)

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
	logProducer kafka.LogProducer
}

func NewLogHandler(lp kafka.LogProducer) LogHandler {
	return &logHandler{logProducer: lp}
}

func (l *logHandler) HandleLogs(c *fiber.Ctx) error {

	var log models.LogEntry

	if err := c.BodyParser(&log); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "failed to parse body")

	}

	byt, err := json.Marshal(log)

	if err != nil {
		return utils.Fail(c, fiber.ErrInternalServerError.Code, "failed to marshal body")
	}

	err = l.logProducer.PushLogsToQueue(byt)

	if err != nil {
		return utils.Fail(c, fiber.ErrInternalServerError.Code, "failed to marshal body")
	}

	return utils.Success(c, "successfully published log.")

}
