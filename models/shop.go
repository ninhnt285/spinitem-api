package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/rs/xid"
	"gopkg.in/mgo.v2/bson"

	"../helpers/config"
)

// Shop is shop model
type Shop struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectId `bson:"user_id" json:"user_id"`
	Platform    string        `bson:"platform" json:"platform"`
	PlatformKey string        `bson:"platform_key" json:"key"`
	IsVerified  bool          `bson:"is_verified" json:"is_verified"`
	IsActive    bool          `bson:"is_active" json:"is_active"`
	Created     time.Time     `bson:"created" json:"created"`
	Products    interface{}
}

// ShopFull includes private params
type ShopFull struct {
	Shop    `bson:",inline"`
	Session string  `bson:"session" json:"session"`
	Shopify Shopify `bson:"shopify" json:"shopify"`
}

const shopCollectionName = "shop"

// Add new shop to DB
// Note 1: Check shop is existed or not
// - If a shop had the same user_id, return that shop
// - If a shop had access_token and different user_id, return error
// - Else, add new shop to DB
// Note 2: Accept update AccessToken when scope changed
func (shop *ShopFull) Add() error {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(shopCollectionName)
	// Make sure no shop with access_token existed before
	var oldShops []ShopFull
	coll.Find(bson.M{"platform": shop.Platform, "platform_key": shop.PlatformKey, "access_token": bson.M{"$exists": true}}).All(&oldShops)
	for _, oldShop := range oldShops {
		if oldShop.UserID == shop.UserID {
			*shop = oldShop
			return nil
		}
		if oldShop.IsVerified {
			return errors.New("The shop was added by other account")
		}
	}
	// Create new session and add shop to DB
	shop.ID = bson.NewObjectId()
	session := md5.Sum([]byte(shop.Platform + shop.PlatformKey + xid.New().String()))
	shop.Session = hex.EncodeToString(session[:])
	shop.Created = time.Now()
	err := coll.Insert(shop)
	if err != nil {
		return err
	}
	return nil
}

// Verify get access_token, load Products, then run callback
func (shop *ShopFull) Verify() error {
	// Get AccessToken
	return nil
}

/*
// UpdateAccess AccessToken of a Shop
func (shop *ShopFull) UpdateAccess(updateShop ShopFull) error {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(shopCollectionName)
	// Compare values
	if updateShop.AccessToken != "" {
		shop.AccessToken = updateShop.AccessToken
	}
	// Update shop
	err := coll.UpdateId(shop.ID, &shop)
	if err != nil {
		return err
	}
	// Delete other same shop has no AccessToken
	_, err = coll.RemoveAll(bson.M{"platform": shop.Platform, "platform_key": shop.PlatformKey, "access_token": bson.M{"$exists": false}})
	return err
}
*/

// GetShopBySession return ShopFull by session
func GetShopBySession(session string) (*ShopFull, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(shopCollectionName)
	// Find Shop
	var shop ShopFull
	err := coll.Find(bson.M{"session": session}).One(&shop)
	return &shop, err
}

// GetShop get a shop by id
func GetShop(id string) (*Shop, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(shopCollectionName)
	// Convert id to bson.ObjectId
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Shop ID is invalid")
	}
	shopID := bson.ObjectIdHex(id)
	// Find shop by ID in DB
	var shop ShopFull
	err := coll.FindId(shopID).One(&shop)
	if err != nil {
		return nil, err
	}
	// Load products of shop
	switch shop.Platform {
	case "shopify":
		// shop.Products, err = ShopifyLoadProducts(shop)
	}
	return &shop.Shop, err
}

// GetAllShops get all shops of user
func GetAllShops(userID string) ([]Shop, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(shopCollectionName)
	// Convert id to bson.ObjectId
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("User ID is invalid")
	}
	bsonUserID := bson.ObjectIdHex(userID)
	var shops []Shop
	err := coll.Find(bson.M{"user_id": bsonUserID}).All(&shops)
	return shops, err
}
