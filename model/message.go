package model

import (
	"chat/config"
	"chat/db"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

const (
	MESSAGE_TYPE_CHAT = iota
	MESSAGE_TYPE_GROUP
)

type Message struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	UserID        bson.ObjectId `json:"user_id" bson:"user_id"`
	RoomID        bson.ObjectId `json:"room_id" bson:"room_id"`
	Type          int           `json:"type" bson:"type"`
	Content       string        `json:"content" bson:"content"`
	CreateTime    int64         `json:"create_time" bson:"create_time"`
	Status        int           `json:"-" bson:"status"`
}

type MessageResult struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	UserID        bson.ObjectId `json:"user_id" bson:"user_id"`
	RoomID        bson.ObjectId `json:"room_id" bson:"room_id"`
	Type          int           `json:"type" bson:"type"`
	Content       string        `json:"content" bson:"content"`
	CreateTime    int64         `json:"create_time" bson:"create_time"`
	UserName      string        `json:"user_name" bson:"user_name"`
	UserAvatar    string        `json:"user_avatar" bson:"user_avatar"`
}

func (message *Message) Update() error {
	if !message.UserID.Valid() || !message.RoomID.Valid() {
		return errors.New("ERROR: object id is not vaild")
	}
	if !message.ID.Valid() {
		message.ID = bson.NewObjectId()
	}
	content := strings.TrimSpace(message.Content)
	if len(content) == 0 {
		return errors.New("ERROR: content can not be empty")
	}
	if message.CreateTime == 0 {
		message.CreateTime = time.Now().Unix()
	}

	session := db.GetMgoSession()
	defer session.Close()
	if message.RoomID.Valid() {
		if err := session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ROOM).UpdateId(message.RoomID, bson.M{"$set": bson.M{"update_time": time.Now().Unix(), "latest_message": message.Content}}); err != nil {
			return err
		}
	}

	_, err := session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_MESSAGE).UpsertId(message.ID, message)
	return err
}

func GetRoomMessageByPage(id string, page, pageSize int) ([]MessageResult, error) {
	objectId := bson.ObjectIdHex(id)
	if !objectId.Valid() {
		return nil, errors.New("ERROR: object id is not valid")
	}

	session := db.GetMgoSession()
	defer session.Close()
	messages := make([]MessageResult, 0)
	aggregate := []bson.M{
		bson.M{
			"$match": bson.M{
				"room_id": objectId,
				"status":  0,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "user",
				"foreignField": "_id",
				"localField":   "user_id",
				"as":           "user",
			},
		},
		bson.M{
			"$unwind": "$user",
		},
		bson.M{
			"$sort": bson.M{
				"create_time": -1,
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":         "$_id",
				"user_id":     "$user_id",
				"room_id":     "$room_id",
				"type":        "$type",
				"content":     "$content",
				"create_time": "$create_time",
				"user_name":   "$user.name",
				"user_avatar": "$user.avatar",
			},
		},
		bson.M{
			"$skip": (page - 1) * pageSize,
		},
		bson.M{
			"$limit": page * pageSize,
		},
	}
	return messages, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_MESSAGE).Pipe(aggregate).All(&messages)
	//return messages, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_MESSAGE).Find(bson.M{"room_id": objectId, "status": 0}).Sort("create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&messages)
}

func GetUserRoomUnReadMessage(roomId, userId string) ([]Message, error) {
	roomIdObject := bson.ObjectIdHex(roomId)
	userIdObject := bson.ObjectIdHex(userId)
	if !roomIdObject.Valid() || !userIdObject.Valid() {
		return nil, errors.New("room id or user id is not vaild")
	}
	session := db.GetMgoSession()
	defer session.Close()

	messages := make([]Message, 0)
	return messages, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_MESSAGE).Find(bson.M{"user_id": userIdObject.Hex(), "room_id": roomIdObject.Hex(), "status": 0}).Sort("create_time").All(&messages)
}

func GetUserRoomUnReadMessageCount(roomId, userId string) (int, error) {
	roomIdObject := bson.ObjectIdHex(roomId)
	userIdObject := bson.ObjectIdHex(userId)
	if !roomIdObject.Valid() || !userIdObject.Valid() {
		return 0, errors.New("room id or user id is not vaild")
	}
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_MESSAGE).Find(bson.M{"user_id": userIdObject, "room_id": roomIdObject, "status": 0}).Count()
}

func SetUserRoomMessageRead(roomId, userId string) error {
	roomIdObject := bson.ObjectIdHex(roomId)
	userIdObject := bson.ObjectIdHex(userId)
	if !roomIdObject.Valid() || !userIdObject.Valid() {
		return errors.New("room id or user id is not vaild")
	}
	session := db.GetMgoSession()
	defer session.Close()

	_, err := session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_MESSAGE).UpdateAll(bson.M{"room_id": roomIdObject, "user_id": userIdObject}, bson.M{"$set": bson.M{"status": 1}})
	return err
}
