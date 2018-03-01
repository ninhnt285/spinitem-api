package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"../../helpers/config"
	"../../models"
)

// User is user model
type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username string        `bson:"username" json:"username"`
	Email    string        `bson:"email" json:"email"`
}

// UserFull includes private params
type UserFull struct {
	User     `bson:",inline"`
	Password string `bson:"password" json:"password"`
}

const collectionName = "user"

// Add new user to DB
// Note: Only new user can add
func (user *UserFull) Add() error {
	dbSession := models.Session.Clone()
	defer dbSession.Close()
	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)

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
	if err != nil {
		return err
	}
	return nil
}

// GetByEmailOrUsername get user by email, password
func GetByEmailOrUsername(email string) (*UserFull, error) {
	dbSession := models.Session.Clone()
	defer dbSession.Close()

	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)

	var user UserFull
	err := coll.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID get user by id
func GetByID(id string) (*User, error) {
	dbSession := models.Session.Clone()
	defer dbSession.Close()

	conf := config.GetInstance()
	coll := dbSession.DB(conf.MongoDatabase).C(collectionName)

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("User ID is invalid")
	}
	userID := bson.ObjectIdHex(id)

	var user User
	err := coll.FindId(userID).One(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// AllUsers get all users
func AllUsers() ([]*User, error) {
	return nil, nil
}
