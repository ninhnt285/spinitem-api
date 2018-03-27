package shopify

import (
	goshopify "github.com/getconversio/go-shopify"

	shopModel "../../models/shop"
	"../config"
)

// Verify verifies for Shopify
type Verify struct {
	Shop      string `json:"shop"`
	State     string `json:"state"`
	Code      string `json:"code"`
	Timestamp string `json:"timestamp"`
	Hmac      string `json:"hmac"`
}

// GetFrom return Verify struct from data
func (v Verify) GetFrom(verifyData map[string]interface{}) (*Verify, error) {
	return nil, nil
}

// Check return AccessCode or Error
func (v *Verify) Check() (*shopModel.ShopFull, error) {
	// Initial
	conf := config.GetInstance()

	// Load shop first
	savedShop, err := shopModel.GetBySession(v.State)
	if err != nil {
		return nil, err
	}

	// Create Shopify App
	app := goshopify.App{
		ApiKey:    conf.ShopifyAPIKey,
		ApiSecret: conf.ShopifyAPISecretKey,
		Scope:     "read_themes,write_themes",
	}

	// Check hmac
	// app.VerifyMessage(message, v.Hmac)

	// Get token
	token, err := app.GetAccessToken(v.Shop, v.Code)
	if err != nil {
		return nil, err
	}

	// Send a test API to Shopify server

	// Update shop with new token
	var updateShop = shopModel.ShopFull{}
	updateShop.AccessCode = token
	err = savedShop.UpdateAccess(updateShop)

	return savedShop, err
}
