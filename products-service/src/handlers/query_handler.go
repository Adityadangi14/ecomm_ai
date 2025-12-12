package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/helpers"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/llm"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
	"github.com/Adityadangi14/ecomm_ai/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
)

type QueryHandler interface {
	GetAiResponse(*fiber.Ctx) error
}

type queryHandler struct {
	aiClient llm.Aiclient
	rdb      *redis.Client
}

func NewQueryHandler(aiClient llm.Aiclient, rdb *redis.Client) QueryHandler {
	return &queryHandler{
		aiClient: aiClient,
		rdb:      rdb,
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

	response := ""

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			select {
			case msg, ok := <-msgChan:
				if !ok {
					// Channel closed, end stream
					fmt.Fprintf(w, "data: [DONE]\n\n")
					w.Flush()
					err := SummerizePastChats(q.rdb, q.aiClient, response, query)
					if err != nil {
						fmt.Println("error summarizing chat ,", err)
					}
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
				response = response + msg.Chunk
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

func SummerizePastChats(rdb *redis.Client, aiClient llm.Aiclient, response string, parms models.AiQueryParams) error {
	key := helpers.GetUserChatKey(parms)
	res, err := helpers.GetUserChat(context.Background(), rdb, key)
	if err != nil {
		slog.Error("error getting user chat", "error", err)
		return err
	}

	chat := map[string]string{
		"query":       parms.Query,
		"ai-response": response,
	}

	str, err := json.Marshal(chat)

	if res == "" {
		err := helpers.SetUserChat(rdb, key, string(str), context.Background())
		slog.Error("Error setting user chat ", "error", err)
		return nil
	}

	slog.Info("summerziation inputs", "response", res, "query", string(str))
	summary := aiClient.SummerizePastChats(res, string(str))

	err = helpers.SetUserChat(rdb, key, summary, context.Background())
	slog.Info("summerize chats", "chat", res)
	if err != nil {
		return err
	}

	return nil

}
