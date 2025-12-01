package repository

import (
	"context"
	"fmt"

	"github.com/Adityadangi14/ecomm_ai/pkg/WDB"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	weaviategraphql "github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

type ProductRepository interface {
	SaveProduct(ctx context.Context, data map[string]any) error
	NearSearchProducts(ctx context.Context, query string, orgID string) ([]map[string]any, error)
	DeleteAllProducts(ctx context.Context) error
}

type prodRepo struct {
	WDB *WDB.WDB
}

func NewProductRepository(wdb *WDB.WDB) ProductRepository {
	return &prodRepo{WDB: wdb}
}

func (p *prodRepo) SaveProduct(ctx context.Context, data map[string]any) error {
	_, err := p.WDB.DB.Data().Creator().WithClassName("Product").WithProperties(data).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *prodRepo) NearSearchProducts(ctx context.Context, query string, orgID string) ([]map[string]any, error) {
	hybrid := p.WDB.DB.GraphQL().HybridArgumentBuilder().
		WithQuery(query).
		WithAlpha(0.8)
	whereFilter := filters.Where().
		WithPath([]string{"orgId"}).
		WithOperator(filters.Equal).
		WithValueText(orgID)

	resp, err := p.WDB.DB.GraphQL().Get().
		WithClassName("Product").
		WithFields(
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_1_associateValue"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_1_associateValueName"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_1_attributeName"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_1_image"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_1_onClickUrl"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_1_price"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_1_skuId"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_1_value"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_2_associateValue"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_2_associateValueName"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_2_attributeName"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_2_image"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_2_onClickUrl"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_2_price"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "attr_2_value"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "brand"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "description"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "name"}),
			weaviategraphql.Field(weaviategraphql.Field{Name: "priceCurrency"}),
		).
		WithLimit(5).
		WithWhere(whereFilter).
		WithHybrid(hybrid).
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to get products %v", err)
	}

	var products []map[string]any

	getMap, ok := resp.Data["Get"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid Get response format")
	}

	fmt.Println("near search products", products)
	rawProducts, ok := getMap["Product"].([]any)
	if !ok {
		return nil, fmt.Errorf("invalid Product result format")
	}

	for _, p := range rawProducts {
		if productMap, ok := p.(map[string]any); ok {
			products = append(products, productMap)
		}
	}

	return products, nil

}

func (p *prodRepo) DeleteAllProducts(ctx context.Context) error {
	err := p.WDB.DB.Schema().
		ClassDeleter().
		WithClassName("Product").
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
