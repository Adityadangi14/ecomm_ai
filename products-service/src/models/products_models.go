package models

type ProductsModel struct {
	Products []struct {
		OrgID         string `json:"orgId"`
		ID            string `json:"id"`
		Name          string `json:"name"`
		Brand         string `json:"brand"`
		Description   string `json:"description"`
		PriceCurrency string `json:"priceCurrency"`
		ProdAttr      []struct {
			SkuID         string `json:"skuId"`
			AttributeName string `json:"attributeName"`
			Value         string `json:"value"`
			Image         string `json:"image"`
			Price         string `json:"price"`
			OnClickURL    string `json:"onClickUrl"`
		} `json:"prodAttr"`
		SearchText string `json:"searchText"`
	} `json:"products"`
}
