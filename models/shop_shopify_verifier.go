package models

// ShopifyVerifier manage verify data of shopify
type ShopifyVerifier struct {
	Shop      string `json:"shop"`
	State     string `json:"state"`
	Code      string `json:"code"`
	Timestamp string `json:"timestamp"`
	Hmac      string `json:"hmac"`
}

func (verifier *ShopifyVerifier) getAccess() (string, error) {
	return "", nil
}
