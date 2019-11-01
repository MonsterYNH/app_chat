package controller

import (
	"chat/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type RoomController struct {}

func (controller *RoomController) GetUserRoom(c *gin.Context) {
	userEntry, _ := c.Get("user")
	user := userEntry.(*model.User)
	rooms, err := model.GetRooms(user.ID.Hex())
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}

	results := make([]model.Room, 0)
	for _, entry := range rooms {
		if entry.Type == model.ROOM_TYPE_FRIEND {
			user1Id := entry.Members[0]
			user2Id := entry.Members[1]
			if user.ID.Hex() != user1Id.Hex() {
				friend, err := model.GetUserById(user1Id.Hex())
				if err != nil {
					model.Result(c, 111, nil, err)
					return
				}
				entry.Title = friend.Name
				entry.Avatar = friend.Avatar
				unRead, err := model.GetUserRoomUnReadMessageCount(entry.ID.Hex(), user1Id.Hex())
				if err != nil {
					model.Result(c, 111, nil, err)
					return
				}
				entry.UnRead = unRead
			}
			if user.ID.Hex() != user2Id.Hex() {
				friend, err := model.GetUserById(user2Id.Hex())
				if err != nil {
					model.Result(c, 111, nil, err)
					return
				}
				entry.Title = friend.Name
				entry.Avatar = friend.Avatar
				unRead, err := model.GetUserRoomUnReadMessageCount(entry.ID.Hex(), user2Id.Hex())
				if err != nil {
					model.Result(c, 111, nil, err)
					return
				}
				entry.UnRead = unRead
			}
		}
		results = append(results, entry)
	}
	model.Result(c, 111, results, nil)
}

func (controller *RoomController) GetRoomUser(c *gin.Context) {
	id := c.Param("id")
	if !bson.ObjectIdHex(id).Valid() {
		model.Result(c, 111, nil, errors.New("ERROR: id is not valid"))
		return
	}
	users, err := model.GetRoomUsers(id)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, users, nil)
}

func (controller *RoomController) GetRoomMessageByPage(c *gin.Context) {
	id := c.Param("id")
	if !bson.ObjectIdHex(id).Valid() {
		model.Result(c, 111, nil, errors.New("ERROR: id is not valid"))
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	messages, err := model.GetRoomMessageByPage(id, page, pageSize)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, messages, nil)
}

func (controller *RoomController) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var update model.Room
	if err := c.BindJSON(&update); err != nil {
		model.Result(nil, 111, nil, err)
		return
	}
	room, err := model.GetRoomById(id)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	if len(update.Title) != 0 && update.Title != room.Title {
		room.Title = update.Title
	}
	if len(update.Description) != 0 && update.Description != room.Description {
		room.Description = update.Description
	}
	if err := room.Update(); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, room, err)
}

func (controller *RoomController) GetFriendRoom(c *gin.Context) {
	id := c.Query("id")
	idObject := bson.ObjectIdHex(id)
	if !idObject.Valid() {
		model.Result(c, 111, nil, errors.New("ERROR: id is not valid"))
		return
	}
	userEntry, _ := c.Get("user")
	user := userEntry.(*model.User)

	friend, err := model.GetUserById(id)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}

	room, err := model.GetRoomByType(model.ROOM_TYPE_FRIEND, []bson.ObjectId{idObject, user.ID})
	if err != nil {
		if err == mgo.ErrNotFound {
			room.Type = model.ROOM_TYPE_FRIEND
			room.Avatar = friend.Avatar
			room.Title = friend.Name
			model.Result(c, 111, room, nil)
			return
		}
		model.Result(c, 111, nil, err)
		return
	}
	room.Avatar = friend.Avatar
	room.Title = friend.Name
	model.Result(c, 111, room, nil)
}
