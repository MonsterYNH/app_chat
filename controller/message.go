package controller

import (
	"chat/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MessageController struct {}

func (controller *MessageController) PostRoomMessage(c *gin.Context) {
	var message model.Message
	if err := c.BindJSON(&message); err != nil {
		model.Result(nil, 111, nil, err)
		return
	}

	if len(message.Content) == 0 {
		model.Result(c, 111, nil, errors.New("ERROR: user id or message content is not vaild"))
		return
	}

	userEntry, _ := c.Get("user")
	createUser := userEntry.(*model.User)

	var room *model.Room
	// 房间不存在先创建房间
	if !message.RoomID.Valid() {
		var err error
		room, err = model.GetRoomByType(model.ROOM_TYPE_FRIEND, []bson.ObjectId{createUser.ID, message.UserID})
		if err != nil && err == mgo.ErrNotFound {
			friend, err := model.GetUserById(c.Query("id"))
			if err != nil {
				model.Result(c, 111, nil, errors.New(fmt.Sprintf("ERROR: get friend failed, error: %s", err)))
				return
			}
			room = &model.Room{
				Type: model.ROOM_TYPE_FRIEND,
				LatestMessage: message.Content,
				Members: []bson.ObjectId{createUser.ID, friend.ID},
				CreateUser: createUser.ID,
			}
			if err := room.Update(); err != nil {
				model.Result(c, 111, nil, err)
				return
			}
		} else if err != nil && err != mgo.ErrNotFound {
			model.Result(c, 111, nil, err)
			return
		}
	} else {
		var err error
		room, err = model.GetRoomById(message.RoomID.Hex())
		if err != nil {
			model.Result(c, 111, nil, err)
			return
		}
		// 判断用户在不在房间
		isExist := false
		for _, entry := range room.Members {
			if createUser.ID.Hex() == entry.Hex() {
				isExist = true
				break
			}
		}
		if !isExist {
			model.Result(c, 111, nil, errors.New("ERROR: user is not in the room"))
			return
		}
	}

	// 创建消息
	message.UserID = createUser.ID
	message.RoomID = room.ID
	message.Type = model.MESSAGE_TYPE_CHAT
	if err := message.Update(); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, message, nil)
}
