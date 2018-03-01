package item

import (
	"errors"

	"gopkg.in/mgo.v2/bson"

	"../../helpers/config"
	"../../models"
)

// Item includes info and all images data
type Item struct {
	ID       bson.ObjectId   `bson:"_id" json:"id"`
	UserID   bson.ObjectId   `bson:"user_id" json:"user_id"`
	Title    string          `bson:"title" json:"title"`
	IsActive bool            `bson:"is_active" json:"is_active"`
	ImageIds []bson.ObjectId `bson:"images" json:"images"`
}

const (
	collectionName = "item"
)

// Add new item to DB
func (item *Item) Add() error {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
	// Validate fields
	if item.ImageIds == nil {
		item.ImageIds = []bson.ObjectId{}
	}
	// Insert item to DB
	item.ID = bson.NewObjectId()
	err := coll.Insert(&item)
	return err
}

// Update an existed item
func (item *Item) Update() error {
	dbSession := models.Session.Clone()
	defer dbSession.Close()

	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
	// Update item
	err := coll.UpdateId(item.ID, &item)
	return err
}

// GetByID get an item by id
func GetByID(id string) (*Item, error) {
	dbSession := models.Session.Clone()
	defer dbSession.Close()

	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
	// Convert id to bson.ObjectId
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Item ID is invalid")
	}
	itemID := bson.ObjectIdHex(id)
	// Find item by ID in DB
	var item Item
	err := coll.FindId(itemID).One(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetAll get all items of user
func GetAll(userID string) ([]Item, error) {
	dbSession := models.Session.Clone()
	defer dbSession.Close()

	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
	// Convert id to bson.ObjectId
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("User ID is invalid")
	}
	bsonUserID := bson.ObjectIdHex(userID)
	var items []Item
	err := coll.Find(bson.M{"user_id": bsonUserID}).All(&items)
	return items, err
}

// Delete item by id
func Delete(id string) error {
	dbSession := models.Session.Clone()
	defer dbSession.Close()

	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
	// Convert id to bson.ObjectId
	if !bson.IsObjectIdHex(id) {
		return errors.New("Item ID is invalid")
	}
	itemID := bson.ObjectIdHex(id)
	// Remove item from DB
	err := coll.RemoveId(itemID)
	if err != nil {
		return err
	}
	return nil
}
