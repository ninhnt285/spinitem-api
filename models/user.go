package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"../helpers/config"
)

// User is user model
type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username string        `bson:"username" json:"username"`
	Email    string        `bson:"email" json:"email"`
	Fullname string        `bson:"fullname" json:"fullname"`
}

// UserFull includes private params
type UserFull struct {
	User     `bson:",inline"`
	Password string `bson:"password" json:"password"`
}

const userCollectionName = "user"

// Add new user to DB
// Note: Only new user can add
func (user *UserFull) Add() error {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(userCollectionName)
	// Make sure user is not existed
	existedUser := UserFull{}
	err := coll.Find(bson.M{"$or": []bson.M{bson.M{"username": user.Username}, bson.M{"email": user.Email}}}).One(&existedUser)
	if err == nil {
		err = errors.New("Email or Username is existed")
		return err
	}
	// Create password by Bcrypt
	newPw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(newPw)
	// Save new user to DB
	user.ID = bson.NewObjectId()
	err = coll.Insert(user)
	return err
}

// GetUserByEmailOrUsername get user by email, password
func GetUserByEmailOrUsername(email string) (*UserFull, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(userCollectionName)
	// Find user by email
	var user UserFull
	err := coll.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUser get user by ID
func GetUser(id string) (*User, error) {
	dbSession := DBSession.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(userCollectionName)
	// Convert id to ObjectID
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("User ID is invalid")
	}
	userID := bson.ObjectIdHex(id)
	// Get User
	var user User
	err := coll.FindId(userID).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers return all users
func GetAllUsers() ([]*User, error) {
	return nil, nil
}
