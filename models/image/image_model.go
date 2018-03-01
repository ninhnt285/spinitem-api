package image

import (
	uf "../../helpers/uploadfile"
	"gopkg.in/mgo.v2/bson"
)

// Image save info of image
type Image struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	File         uf.UploadFile `bson:",inline"`
	Index        int           `bson:"index" json:"index"`
	CaptureIndex int           `bson:"capture_index" json:"capture_index"`
	ItemID       bson.ObjectId `bson:"item_id" json:"item_id"`
}

// Add new image to Database
func (image *Image) Add() error {
	return nil
}
