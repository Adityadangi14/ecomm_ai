package schema

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func CreateProductClass(client *weaviate.Client) error {
	ctx := context.Background()

	// Check if class already exists
	existing, err := client.Schema().Getter().Do(ctx)
	if err != nil {
		panic(err)
	}

	for _, c := range existing.Classes {
		if c.Class == "Product" {
			fmt.Println("Class Product already exists, skipping creation")
			return nil
		}
	}

	// Define new schema
	productClass := &models.Class{
		Class:           "Product",
		VectorIndexType: "hnsw",
		Vectorizer:      "text2vec-transformers",
		ModuleConfig: map[string]interface{}{
			"text2vec-transformers": map[string]any{
				"vectorizeClassName": false,
			},
		},
		Properties: []*models.Property{
			{
				Name:     "search_text",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": false, // vectorized
					},
				},
			},
			{
				Name:     "name",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "brand",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "description",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "product_image_description",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_1_associateValue",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_1_associateValueName",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_1_value",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_1_attributeName",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_1_image",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_1_onClickUrl",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},

			{
				Name:     "attr_1_skuId",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_2_associateValue",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_2_associateValueName",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_2_attributeName",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_2_image",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_2_onClickUrl",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_2_price",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_2_skuId",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "orgId",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "productId",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "priceCurrency",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_1_price",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name:     "attr_2_value",
				DataType: []string{"text"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": true,
					},
				},
			},
		},
	}

	return client.Schema().ClassCreator().WithClass(productClass).Do(ctx)
}
