package llm

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
	"github.com/Adityadangi14/ecomm_ai/utils"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Aiclient interface {
	ImageByteToText(context.Context, []byte) (string, error)
	SummerizePastQueris([]string) string
	GetSementicText(map[string]any) string
}

type aiclient struct {
	LlmClient *openai.Client
}

func NewAiClient() Aiclient {
	key := os.Getenv("OPENAI_KEY")
	client := openai.NewClient(option.WithAPIKey(key))
	return &aiclient{LlmClient: &client}
}

func (a *aiclient) ImageByteToText(ctx context.Context, imgBytes []byte) (string, error) {
	// Convert image bytes to base64
	b64 := base64.StdEncoding.EncodeToString(imgBytes)
	imageURLData := "data:image/jpeg;base64," + b64

	// Construct chat completion parameters
	params := openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4o, // adjust as per your model name
		Messages: []openai.ChatCompletionMessageParamUnion{
			// User message with image and instruction
			openai.UserMessage(fmt.Sprintf(
				"Here is a product image: [image: %s]\n\n"+
					"Please provide a detailed ecommerce-style description of the product in the image: material, color, design, brand cues, use-case, features.",
				imageURLData,
			)),
		},
		Temperature: openai.Float(0.7),
	}

	// Send request
	resp, err := a.LlmClient.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("openai chat completion error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func (a *aiclient) SummerizePastQueris([]string) string {
	return ""
}

func (a *aiclient) GetSementicText(map[string]any) string {
	return ""
}

func (a *aiclient) ProcessProduct(prod models.Product) error {
	prodMap := prod.ToFlatMap()

	byt := utils.GetImageBytesFromFlatMap(prodMap)
	ctx := context.Background()
	res, err := a.ImageByteToText(ctx, byt)

	if err != nil {
		return err
	}

	prodMap["productImageDiscription"] = res

	return nil

}
