package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
	"github.com/Adityadangi14/ecomm_ai/utils"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Aiclient interface {
	ImageByteToText(context.Context, string) (string, error)
	SummerizePastQueris([]string) string
	GetSementicText(map[string]any) (string, error)
	ProcessProduct(prod models.Product) (map[string]any, error)
	GetAiQueryReponse(params models.AiQueryParams, msg chan models.MessageChanStruct)
}

type aiclient struct {
	LlmClient *openai.Client
}

func NewAiClient() Aiclient {
	key := os.Getenv("OPENAI_KEY")
	client := openai.NewClient(option.WithAPIKey(key))
	return &aiclient{LlmClient: &client}
}

func (a *aiclient) ImageByteToText(ctx context.Context, url string) (string, error) {

	params := openai.ChatCompletionNewParams{
		Model: openai.ChatModel("gpt-4o-mini"),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{

				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{

						OfArrayOfContentParts: []openai.ChatCompletionContentPartUnionParam{
							openai.TextContentPart("Please describe this product image: material, color, design, brand cues, and use-case."),
							openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{URL: url}),
						},
					},
				},
			},
		},
		Temperature: openai.Float(0.7),
	}

	var resp *openai.ChatCompletion
	var err error

	// Retry logic
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = a.LlmClient.Chat.Completions.New(ctx, params)
		if err == nil {
			break
		}

		// Exponential backoff sleep
		wait := time.Duration(attempt*attempt) * time.Second
		log.Printf("Retrying OpenAI request after error (attempt %d/%d): %v. Waiting %v...", attempt, maxRetries, err, wait)
		time.Sleep(wait)
	}

	if err != nil {
		return "", fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

func (a *aiclient) SummerizePastQueris([]string) string {
	return ""
}

func (a *aiclient) GetSementicText(prod map[string]any) (string, error) {
	ctx := context.Background()

	// Convert product map to JSON string
	jsonBytes, err := json.Marshal(prod)
	if err != nil {
		return "", fmt.Errorf("failed to marshal product map: %w", err)
	}

	// Build prompt
	finalPrompt := fmt.Sprintf("%s\n\n%s\n\nJSON:\n%s",
		SEMENTIC_SEARCH_PRODUCT_PROMPT,
		SEMENTIC_SEARCH_PRODUCT_OUTPUT_PROMPT,
		string(jsonBytes),
	)

	// Construct chat completion parameters
	params := openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(finalPrompt),
		},
		Temperature: openai.Float(0.4),
	}

	var resp *openai.ChatCompletion
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = a.LlmClient.Chat.Completions.New(ctx, params)
		if err == nil {
			break
		}

		// Backoff strategy: 1s, 2s, 4s (or use square: 1,4,9)
		wait := time.Duration(attempt*attempt) * time.Second
		log.Printf("Retrying semantic text request (attempt %d/%d): %v. Waiting %v...",
			attempt, maxRetries, err, wait,
		)
		time.Sleep(wait)
	}

	if err != nil {
		return "", fmt.Errorf("GetSementicText failed after %d attempts: %w", maxRetries, err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func (a *aiclient) ProcessProduct(prod models.Product) (map[string]any, error) {
	ctx := context.Background()
	// Convert to flat map
	prodMap := prod.ToFlatMap()

	// Extract image bytes
	imgUrl := utils.ExtractImageUrlFromFlatMap(prodMap)
	if len(imgUrl) == 0 {
		return nil, fmt.Errorf("no image bytes found for product")
	}

	// Step 1: Extract product image description via AI
	imageDesc, err := a.ImageByteToText(ctx, imgUrl)
	if err != nil {
		return nil, fmt.Errorf("image description generation failed: %w", err)
	}

	prodMap["product_image_description"] = imageDesc

	semanticText, err := a.GetSementicText(prodMap)
	if err != nil {

		return nil, fmt.Errorf("semantic text generation failed: %w", err)
	}

	prodMap["search_text"] = semanticText

	return prodMap, nil
}

func (a *aiclient) GetAiQueryReponse(params models.AiQueryParams, msg chan models.MessageChanStruct) {

	stream := a.LlmClient.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Hello, please stream a response."),
		},
		Temperature: openai.Float(0.7),
	})

	defer stream.Close()
	defer close(msg)
	defer wg.Done()

	if stream.Err() != nil {
		fmt.Println("Stream error:", stream.Err())
	}

	for stream.Next() {
		event := stream.Current()
		if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
			msg <- models.MessageChanStruct{Chunk: event.Choices[0].Delta.Content}
		}
	}

	if err := stream.Err(); err != nil {
		msg <- models.MessageChanStruct{Err: fmt.Errorf("%s", fmt.Sprintf("Stream error: %v", err))}
	}
}
