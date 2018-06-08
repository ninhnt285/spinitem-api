package models

// Shopify save all info of a shopify
type Shopify struct {
	AccessToken string `bson:"access_token" json:"access_token"`
	ShopURL     string `bson:"shop_url" json:"shop_url"`
}

// Callback do callback for Shopify:
// - Get AccessToken
// - Update Shop: shopify, is_verified, is_active
func (shopify *Shopify) Callback(verifier ShopifyVerifier, shop ShopFull) error {
	// Get AccessToken

	// Update ShopFull

	return nil
}
