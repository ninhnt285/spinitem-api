package models

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"../helpers/config"
)

const (
	imageCollectionName = "image"
)

// Image save info of image
type Image struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	ItemID       bson.ObjectId `bson:"item_id" json:"item_id"`
	Index        int           `bson:"index" json:"index"`
	CaptureIndex int           `bson:"capture_index" json:"capture_index"`
	IsActive     bool          `bson:"is_active" json:"is_active"`
	Pitch        float64       `bson:"pitch" json:"pitch"`
	Roll         float64       `bson:"roll" json:"roll"`
	Yaw          float64       `bson:"yaw" json:"yaw"`
	UploadFile   `bson:",inline" json:",inline"`
}

// PrepareResult adds static URL to destination
func (image *Image) PrepareResult() {
	conf := config.GetInstance()
	image.Destination = conf.StaticURL + image.Destination
}

// Add new image to Database
func (image *Image) Add() error {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(imageCollectionName)
	// Create new ObjectID
	image.ID = bson.NewObjectId()
	// Add new Image to DB
	err := coll.Insert(image)
	if err != nil {
		return err
	}
	return nil
}

// Delete image from server and DB
func (image *Image) Delete() error {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(imageCollectionName)
	// Remove files in server
	os.Remove(conf.PublicDir + fmt.Sprintf(image.Destination, ""))
	os.Remove(conf.PublicDir + fmt.Sprintf(image.Destination, "_origin"))
	for _, option := range ImgOptions {
		suffix := "_" + strconv.Itoa(option.Width) + "x" + strconv.Itoa(option.Height)
		if option.Crop {
			suffix += "_square"
		}
		os.Remove(conf.PublicDir + fmt.Sprintf(image.Destination, suffix))
	}
	// Delete Image from DB
	return coll.RemoveId(image.ID)
}

// GetImageByObjectID return image
func GetImageByObjectID(imageID bson.ObjectId) (*Image, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(imageCollectionName)
	// Find image by ID in DB
	var image Image
	err := coll.FindId(imageID).One(&image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImage by id string
func GetImage(id string) (*Image, error) {
	// Convert id to bson.ObjectId
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Image ID is invalid")
	}
	imageID := bson.ObjectIdHex(id)
	return GetImageByObjectID(imageID)
}

// GetAllImages return all images in an item
func GetAllImages(itemID string) ([]Image, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(imageCollectionName)
	// Convert itemID to bson.ObjectId
	if !bson.IsObjectIdHex(itemID) {
		return nil, errors.New("Item ID is invalid")
	}
	bsonItemID := bson.ObjectIdHex(itemID)
	var images []Image
	err := coll.Find(bson.M{"item_id": bsonItemID}).All(&images)
	return images, err
}
