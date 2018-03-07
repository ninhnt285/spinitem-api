package image

import (
	"errors"
	"os"

	"gopkg.in/mgo.v2/bson"

	"../../helpers/config"
	"../../models"
	uf "../uploadfile"
)

const (
	collectionName = "image"
)

// Image save info of image
type Image struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	ItemID        bson.ObjectId `bson:"item_id" json:"item_id"`
	Index         int           `bson:"index" json:"index"`
	CaptureIndex  int           `bson:"capture_index" json:"capture_index"`
	IsActive      bool          `bson:"is_active" json:"is_active"`
	uf.UploadFile `bson:",inline" json:",inline"`
}

// PrepareResult add static URL to destination
func (image *Image) PrepareResult() {
	conf := config.GetInstance()
	image.Destination = conf.StaticURL + image.Destination
}

// Add new image to Database
func (image *Image) Add() error {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
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
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
	// TODO: Remove file in server
	err := os.Remove(conf.PublicDir + image.Destination)
	// Delete Image from DB
	err = coll.RemoveId(image.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetImageByObjectID return image
func GetImageByObjectID(imageID bson.ObjectId) (*Image, error) {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
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

// GetAllImages get all images
func GetAllImages(itemID string) ([]Image, error) {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)
	// Convert itemID to bson.ObjectId
	if !bson.IsObjectIdHex(itemID) {
		return nil, errors.New("Item ID is invalid")
	}
	bsonItemID := bson.ObjectIdHex(itemID)
	var images []Image
	err := coll.Find(bson.M{"item_id": bsonItemID}).All(&images)
	return images, err
}
