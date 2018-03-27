package shop

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/rs/xid"
	"gopkg.in/mgo.v2/bson"

	"../../helpers/config"
	"../../models"
)

// Shop is shop model
type Shop struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserID   bson.ObjectId `bson:"user_id" json:"user_id"`
	Platform string        `bson:"platform" json:"platform"`
	Key      string        `bson:"key" json:"key"`
	Session  string        `bson:"session" json:"session"`
	IsActive bool          `bson:"is_active" json:"is_active"`
	Created  time.Time     `bson:"created" json:"created"`
}

// ShopFull includes private params
type ShopFull struct {
	Shop       `bson:",inline"`
	AccessCode string `bson:"access_code" json:"access_code"`
}

const collectionName = "shop"

// Add new shop to DB
// Note: Check existed shop
// - If a shop had access_code, return error
// - If a shop had the same user_id, return that shop
func (shop *ShopFull) Add() error {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)

	// Make sure no shop with access_code existed
	var oldShops []ShopFull
	coll.Find(bson.M{"platform": shop.Platform, "key": shop.Key, "access_code": bson.M{"$exists": true}}).All(&oldShops)
	for _, oldShop := range oldShops {
		if oldShop.AccessCode != "" {
			return errors.New("The shop was added by other account")
		}
		if (oldShop.Platform == shop.Platform) && (oldShop.Key == shop.Key) {
			*shop = oldShop
			return nil
		}
	}

	// Create new session and add shop to DB
	shop.ID = bson.NewObjectId()
	session := md5.Sum([]byte(shop.Platform + shop.Key + xid.New().String()))
	shop.Session = hex.EncodeToString(session[:])
	shop.Created = time.Now()
	err := coll.Insert(shop)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAccess AccessCode of a Shop
func (shop *ShopFull) UpdateAccess(updateShop ShopFull) error {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)

	// Compare values
	if updateShop.AccessCode != "" {
		shop.AccessCode = updateShop.AccessCode
	}
	// Update shop
	err := coll.UpdateId(shop.ID, &shop)
	if err != nil {
		return err
	}
	// Delete other same shop has no AccessCode
	_, err = coll.RemoveAll(bson.M{"platform": shop.Platform, "key": shop.Key, "access_code": bson.M{"$exists": false}})

	return err
}

// GetBySession return ShopFull by session
func GetBySession(session string) (*ShopFull, error) {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)

	// Find Shop
	var shop ShopFull
	err := coll.Find(bson.M{"session": session}).One(&shop)

	return &shop, err
}

// GetAll get all shops of user
func GetAll(userID string) ([]Shop, error) {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
	// Convert id to bson.ObjectId
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("User ID is invalid")
	}
	bsonUserID := bson.ObjectIdHex(userID)
	var shops []Shop
	err := coll.Find(bson.M{"user_id": bsonUserID}).All(&shops)
	return shops, err
}
