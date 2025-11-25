package llm

const (
	SEMENTIC_SEARCH_PRODUCT_PROMPT = `You are an AI assistant that converts structured product JSON data into a rich, human-readable semantic description that captures all important information useful for semantic search, product recommendations, and natural-language retrieval.

Generate a clean paragraph description combining all relevant attributes below. Do not include field names, special characters, JSON formatting, or HTML tags. Include product name, brand, product type, attributes (like size, color), pricing, main features, material or fabric, use case, and any additional context. Write naturally as if describing the product to a customer.`

	SEMENTIC_SEARCH_PRODUCT_OUTPUT_PROMPT = `
	Output:
	- A single semantic searchable text summary
	- No HTML tags or JSON field names
	- No URLs or SKU details
	`
)
