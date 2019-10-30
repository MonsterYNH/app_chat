package model

import (
	"chat/config"
	"chat/db"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	ROOM_TYPE_FRIEND = iota + 1
	ROOM_TYPE_GROUP
)

type Room struct {
	ID            bson.ObjectId   `json:"id" bson:"_id"`
	Type          int             `json:"type" bson:"type"`
	Avatar        string          `json:"avatar" bson:"avatar"`
	Title         string          `json:"title" bson:"title"`
	LatestMessage string          `json:"latest_message" bson:"latest_message"`
	Description   string          `json:"description" bson:"description"`
	Members       []bson.ObjectId `json:"members" bson:"members"`
	CreateTime    int64           `json:"create_time" bson:"create_time"`
	UpdateTime    int64           `json:"update_time" bson:"update_time"`
	CreateUser    bson.ObjectId   `json:"create_user" bson:"create_user"`
	UnRead        int             `json:"un_read" bson:"-"`
	Status        int             `json:"status" bson:"status"`
}

func (room *Room) Update() error {
	session := db.GetMgoSession()
	defer session.Close()

	if !room.CreateUser.Valid() {
		return errors.New("ERROR: create_user can not empty")
	}
	if room.Members == nil || len(room.Members) == 0 {
		return errors.New("ERROR: members can not empty")
	}

	hasOwner := false
	for _, entry := range room.Members {
		if entry.Hex() == room.CreateUser.Hex() {
			hasOwner = true
			break
		}
	}
	if !hasOwner {
		return errors.New("ERROR: members must have creater")
	}
	if len(room.Members) == 2 {
		room.Type = ROOM_TYPE_FRIEND
	} else if len(room.Members) > 2 {
		room.Type = ROOM_TYPE_GROUP
	} else {
		return errors.New("ERROR: wrong room type")
	}

	if !room.ID.Valid() {
		room.ID = bson.NewObjectId()
	}
	if room.CreateTime == 0 {
		room.CreateTime = time.Now().Unix()
	}
	if room.UpdateTime == 0 {
		room.UpdateTime = time.Now().Unix()
	}
	_, err := session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ROOM).UpsertId(room.ID, room)
	return err
}

func GetRooms(id string) ([]Room, error) {
	objectId := bson.ObjectIdHex(id)
	if !objectId.Valid() {
		return nil, errors.New("ERROR: object id is not vaild")
	}
	session := db.GetMgoSession()
	defer session.Close()

	rooms := make([]Room, 0)
	err := session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ROOM).Find(bson.M{
		"status": 0,
		"members": bson.M{
			"$in": []bson.ObjectId{objectId},
		},
	}).Select(bson.M{"status": 0}).Sort("update_time").All(&rooms)
	return rooms, err
}

func GetRoomByUsers(users []bson.ObjectId) (*Room, error) {
	session := db.GetMgoSession()
	defer session.Close()

	room := Room{}
	return &room, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ROOM).Find(bson.M{"members": bson.M{"$in": users}, "type": 1}).One(&room)
}

func GetRoomUsers(id string) ([]User, error) {
	objectId := bson.ObjectIdHex(id)
	if !objectId.Valid() {
		return nil, errors.New("ERROR: object id not vaild")
	}
	session := db.GetMgoSession()
	defer session.Close()

	room := Room{}
	if err := session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ROOM).FindId(objectId).One(&room); err != nil {
		return nil, err
	}

	users := make([]User, 0)
	return users, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_USER).Find(bson.M{"_id": bson.M{"$in": room.Members}}).All(&users)
}

func GetRoomById(id string) (*Room, error) {
	objectId := bson.ObjectIdHex(id)
	if !objectId.Valid() {
		return nil, errors.New("ERROR: object id not vaild")
	}
	session := db.GetMgoSession()
	defer session.Close()

	var room Room
	return &room, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ROOM).FindId(objectId).One(&room)
}

func GetRoomByType(kind int, users []bson.ObjectId) (*Room, error) {
	session := db.GetMgoSession()
	defer session.Close()

	var room Room
	return &room, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ROOM).Find(bson.M{"type": kind, "members": bson.M{"$in": users}}).One(&room)
}
