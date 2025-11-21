package models

import "fmt"

type ProductsModel struct {
	Products []Product `json:"products"`
}

type Product struct {
	OrgID         string     `json:"orgId"`
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Brand         string     `json:"brand"`
	Description   string     `json:"description"`
	PriceCurrency string     `json:"priceCurrency"`
	Attributes    []ProdAttr `json:"prodAttr"`
}

type ProdAttr struct {
	SkuID              string `json:"skuId"`
	AttributeName      string `json:"attributeName"`
	Value              string `json:"value"`
	AssociateValueName string `json:"associateValueName"`
	AssociateValue     string `json:"associateValue"`
	Image              string `json:"image"`
	Price              string `json:"price"`
	OnClickURL         string `json:"onClickUrl"`
}

func (p Product) ToFlatMap() map[string]interface{} {
	m := map[string]interface{}{
		"orgId":         p.OrgID,
		"productId":     p.ID,
		"name":          p.Name,
		"brand":         p.Brand,
		"description":   p.Description,
		"priceCurrency": p.PriceCurrency,
	}

	// Flatten all attributes
	for i, a := range p.Attributes {
		prefix := fmt.Sprintf("attr_%d_", i+1)

		m[prefix+"skuId"] = a.SkuID
		m[prefix+"attributeName"] = a.AttributeName
		m[prefix+"value"] = a.Value
		m[prefix+"associateValueName"] = a.AssociateValueName
		m[prefix+"associateValue"] = a.AssociateValue
		m[prefix+"image"] = a.Image
		m[prefix+"price"] = a.Price
		m[prefix+"onClickUrl"] = a.OnClickURL
	}

	return m
}
