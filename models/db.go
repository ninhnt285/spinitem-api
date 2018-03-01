package models

import (
	"log"

	"../helpers/config"
	mgo "gopkg.in/mgo.v2"
)

// Session save session of Mongo
var Session *mgo.Session

// ConnectDB connect App to MongoDB Service
func ConnectDB(dataSourceName string) {
	cf := config.GetInstance()
	var err error
	Session, err = mgo.Dial(cf.MongoServer)
	if err != nil {
		log.Panic(err)
	}
}

// CloseDB close DB
func CloseDB() {
	Session.Close()
}
