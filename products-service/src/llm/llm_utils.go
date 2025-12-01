package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/helpers"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/repository"
	"github.com/Adityadangi14/ecomm_ai/utils"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/redis/go-redis/v9"
)

type Aiclient interface {
	ImageByteToText(context.Context, string) (string, error)
	SummerizePastQueris(string) string
	GetSementicText(map[string]any) (string, error)
	ProcessProduct(prod models.Product) (map[string]any, error)
	GetAiQueryReponse(params models.AiQueryParams, msg chan models.MessageChanStruct)
	SummerizePastChats(pastSummary string, query string) string
}

type aiclient struct {
	LlmClient   *openai.Client
	rbd         *redis.Client
	productRepo repository.ProductRepository
}

func NewAiClient(rdb *redis.Client, productRepo repository.ProductRepository) Aiclient {
	key := os.Getenv("OPENAI_KEY")
	client := openai.NewClient(option.WithAPIKey(key))
	return &aiclient{LlmClient: &client, rbd: rdb, productRepo: productRepo}
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

func (a *aiclient) SummerizePastQueris(query string) string {

	resp, err := a.LlmClient.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: openai.ChatModelGPT4Turbo,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("Using only the information provided, condense the user's decayed search queries into a single, concise sentence that captures the overall intent and topics, strictly avoiding opinions, assumptions, extra details, or creative additions. Preserve the meaning according to the weight (higher-weight queries influence the summary more), and produce only one clear summary line as output."),
				openai.UserMessage(query),
			},
			MaxTokens: openai.Int(60),
		},
	)

	if err != nil {
		log.Println("OpenAI summary error:", err)
		return ""
	}

	if len(resp.Choices) == 0 {
		return ""
	}

	return resp.Choices[0].Message.Content
}

func (a *aiclient) SummerizePastChats(pastSummary string, query string) string {

	resp, err := a.LlmClient.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: openai.ChatModelGPT4Turbo,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Previous Summary: " + pastSummary),

				// New message to incorporate into summary
				openai.UserMessage("Latest Message: " + query),

				// Explicit instruction to merge them
				openai.UserMessage(CHAT_SUMMARY_PROMPT),
			},
			Temperature: openai.Float(0.2),
		},
	)

	if err != nil {
		log.Println("OpenAI summary error:", err)
		return ""
	}

	if len(resp.Choices) == 0 {
		return ""
	}

	return resp.Choices[0].Message.Content
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

	key := helpers.GetUserQueriesKey(params)

	err := helpers.SetUserQueries(a.rbd, params.Query, key, context.Background())

	if err != nil {
		fmt.Println(err)
	}

	res, err := helpers.GetQueriesWithDecay(context.Background(), a.rbd, key)

	querySummary := a.SummerizePastQueris(res)

	fmt.Println("querysummary", querySummary)

	if err != nil {
		fmt.Println(err)
	}

	products, err := a.productRepo.NearSearchProducts(context.Background(), querySummary, params.OrgID)

	fmt.Println("products", products)

	if err != nil {
		fmt.Println("unable to do near search", err)
	}

	byt, err := json.Marshal(products)

	chatRes, err := helpers.GetUserChat(context.Background(), a.rbd, helpers.GetUserChatKey(params))

	fmt.Println("chatSummary", chatRes)

	stream := a.LlmClient.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4_1,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(fmt.Sprintf("These are the products you can recommend \n %v", string(byt))),
			openai.DeveloperMessage(RESPONSE_UI_COMPONENTS_AND_PROMPT),
			openai.AssistantMessage(fmt.Sprintf("This is previous chat summary \n %v", chatRes)),
			openai.UserMessage(params.Query),
		},
	})

	defer stream.Close()
	defer close(msg)

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
