package models

import (
	"errors"
	"reflect"

	"gopkg.in/mgo.v2/bson"

	"../helpers/config"
)

// Item includes info and all images data
type Item struct {
	ID            bson.ObjectId   `bson:"_id" json:"id"`
	UserID        bson.ObjectId   `bson:"user_id" json:"user_id"`
	Title         string          `bson:"title" json:"title"`
	IsActive      bool            `bson:"is_active" json:"is_active"`
	BackgroundURL string          `bson:"background_url" json:"background_url"`
	ImageIds      []bson.ObjectId `bson:"images" json:"images"`
}

const (
	itemCollectionName = "item"
)

// Add new item to DB
func (item *Item) Add() error {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(itemCollectionName)
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
func (item *Item) Update(updateItem Item) error {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(itemCollectionName)
	// Compare values
	if updateItem.Title != "" {
		item.Title = updateItem.Title
	}
	if updateItem.IsActive != item.IsActive {
		item.IsActive = updateItem.IsActive
	}
	if updateItem.BackgroundURL != item.BackgroundURL {
		item.BackgroundURL = updateItem.BackgroundURL
	}
	if (updateItem.ImageIds != nil) && (!reflect.DeepEqual(item.ImageIds, updateItem.ImageIds)) {
		item.ImageIds = updateItem.ImageIds
	}
	// Update item
	err := coll.UpdateId(item.ID, &item)
	return err
}

// Delete item by id
func (item *Item) Delete() error {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(itemCollectionName)
	// Remove all images first
	images, err := GetAllImages(item.ID.Hex())
	for _, image := range images {
		image.Delete()
	}
	// Remove item from DB
	err = coll.RemoveId(item.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetItem get an item by id
func GetItem(id string) (*Item, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(itemCollectionName)
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

// GetAllItems get all items of user
func GetAllItems(userID string) ([]Item, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(itemCollectionName)
	// Convert id to bson.ObjectId
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("User ID is invalid")
	}
	bsonUserID := bson.ObjectIdHex(userID)
	var items []Item
	err := coll.Find(bson.M{"user_id": bsonUserID}).All(&items)
	return items, err
}
