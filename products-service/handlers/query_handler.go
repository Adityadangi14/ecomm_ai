package handlers

import (
	"bufio"
	"fmt"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/llm"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
	"github.com/Adityadangi14/ecomm_ai/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type QueryHandler interface {
	GetAiResponse(*fiber.Ctx) error
}

type queryHandler struct {
	aiClient llm.Aiclient
}

func NewQueryHandler(aiClient llm.Aiclient) QueryHandler {
	return &queryHandler{
		aiClient: aiClient,
	}
}

func (q *queryHandler) GetAiResponse(c *fiber.Ctx) error {

	var query models.AiQueryParams

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	if err := c.BodyParser(&query); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, fmt.Sprintf("failed to parse body: %v", err))
	}

	msgChan := make(chan models.MessageChanStruct)

	go q.aiClient.GetAiQueryReponse(models.AiQueryParams{Query: query.Query, SessionID: query.SessionID, UserID: query.UserID, OrgID: query.OrgID}, msgChan)

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			select {
			case msg, ok := <-msgChan:
				if !ok {
					// Channel closed, end stream
					fmt.Fprintf(w, "data: [DONE]\n\n")
					w.Flush()
					return
				}

				if msg.Err != nil {
					// Send error as SSE
					fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", msg.Err.Error())
					w.Flush()
					return
				}

				// Send chunk as SSE
				fmt.Fprintf(w, "data: %s\n\n", msg.Chunk)

				err := w.Flush()
				if err != nil {
					// Connection closed by client
					fmt.Printf("Error flushing: %v. Closing connection.\n", err)
					return
				}
			}
		}
	}))

	return nil

}
