package models

import (
	"log"

	"../helpers/config"
	mgo "gopkg.in/mgo.v2"
)

// DBSession save session of Mongo
var DBSession *mgo.Session

// ConnectDB connect App to MongoDB Service
func ConnectDB(dataSourceName string) {
	cf := config.GetInstance()
	var err error
	DBSession, err = mgo.Dial(cf.MongoServer)
	if err != nil {
		log.Panic(err)
	}
}

// CloseDB close DB
func CloseDB() {
	DBSession.Close()
}
