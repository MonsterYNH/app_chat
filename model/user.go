package model

import (
	"chat/config"
	"chat/db"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	ID         bson.ObjectId   `json:"id" bson:"_id"`
	Avatar     string          `json:"avatar" bson:"avatar"`
	Name       string          `json:"name" bson:"name"`
	Signture   string          `json:"signture" bson:"signture"`
	Sex        int             `json:"sex" bson:"sex"`
	Age        int             `json:"age" bson:"age"`
	Account    string          `json:"account" bson:"account"`
	Password   string          `json:"password" bson:"password"`
	CreateTime int64           `json:"create_time" bson:"create_time"`
	UpdateTime int64           `json:"update_time" bson:"update_time"`
	Friends    []bson.ObjectId `json:"friends" bson:"friends"`
	Status     int             `json:"-" bson:"status"`
}

type UserInfo struct {
	ID       *bson.ObjectId `json:"id" bson:"_id"`
	Name     string         `json:"name" bson:"name"`
	Avatar   string         `json:"avatar" bson:"avatar"`
	Signture string         `json:"signture" bson:"signture"`
}

func (user *User) Update() error {
	session := db.GetMgoSession()
	defer session.Close()

	if !user.ID.Valid() {
		user.ID = bson.NewObjectId()
	}
	if user.CreateTime == 0 {
		user.CreateTime = time.Now().Unix()
	}
	if user.UpdateTime == 0 || user.UpdateTime < time.Now().Unix() {
		user.UpdateTime = time.Now().Unix()
	}
	_, err := session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_USER).UpsertId(user.ID, user)
	return err
}

func GetUserById(id string) (*User, error) {
	objectId := bson.ObjectIdHex(id)
	if !objectId.Valid() {
		return nil, errors.New("ERROR: object id is not valid")
	}
	user := User{}
	session := db.GetMgoSession()
	defer session.Close()

	return &user, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_USER).FindId(objectId).One(&user)
}

func GetUserByAccountAndPassword(account, password string) (*User, error) {
	user := User{}
	session := db.GetMgoSession()
	defer session.Close()

	return &user, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_USER).Find(bson.M{"account": account, "password": password}).One(&user)
}

func (user *User) GetUserFriends() ([]UserInfo, error) {
	session := db.GetMgoSession()
	defer session.Close()

	users := make([]UserInfo, 0)
	return users, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_USER).Find(bson.M{"_id": bson.M{"$in": user.Friends}}).Select(bson.M{"friends": 0}).All(&users)
}
