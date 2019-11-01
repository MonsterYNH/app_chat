package main

import (
	"chat/config"
	"chat/controller"
	"chat/db"
	"chat/middleware"
	"chat/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	//createTestUser()
	engineServer := gin.Default()

	var userGroupController = controller.UserController{}
	engineServer.POST("/login", userGroupController.Login)
	engineServer.POST("/regist", userGroupController.Regist)

	v1Api := engineServer.Group("/v1", middleware.AuthMiddleware)

	// user
	userGroup := v1Api.Group("/user")
	{
		userGroup.GET("/list", userGroupController.GetUserRelationShip)
		userGroup.PUT("/update", userGroupController.UpdateUser)
	}

	// room
	roomGroup := v1Api.Group("/room")
	var rommGroupController = controller.RoomController{}
	{
		roomGroup.GET("/list", rommGroupController.GetUserRoom)
		roomGroup.GET("/user/list/:id", rommGroupController.GetRoomUser)
		roomGroup.GET("/message/list/:id", rommGroupController.GetRoomMessageByPage)
		roomGroup.GET("/friend", rommGroupController.GetFriendRoom)
	}

	// message
	messageGroup := v1Api.Group("/message")
	var messageGroupController = controller.MessageController{}
	{
		messageGroup.POST("/create", messageGroupController.PostRoomMessage, middleware.SendUserRoomMiddleware)
		messageGroup.POST("/read", messageGroupController.SetUserRoomMessageRead)
	}

	// article
	articleGroup := v1Api.Group("/article")
	var articleGroupController = controller.ArticleController{}
	{
		articleGroup.POST("/create", articleGroupController.PostCreateArticle)
		articleGroup.GET("/list", articleGroupController.GetArticleList)
		articleGroup.POST("/comment/:id", articleGroupController.PostCreateArticleComment)
	}
	engineServer.Run(config.ENV_SERVER_URL)
}

func createTestUser() {
	db.FlushDb()
	userCuteId := bson.NewObjectId()
	userMonsterId := bson.NewObjectId()
	user := model.User{
		ID: userCuteId,
		Avatar: "http://b-ssl.duitang.com/uploads/item/201410/20/20141020224133_Ur54c.jpeg",
		Name: "你的小可爱",
		Account: "cute",
		Password: "cute",
		Friends: []bson.ObjectId{
			userMonsterId,
		},
	}
	monster := model.User{
		ID: userMonsterId,
		Avatar: "http://b-ssl.duitang.com/uploads/item/201702/05/20170205222154_WLdJS.jpeg",
		Name: "你的小怪兽",
		Account: "monster",
		Password: "monster",
		Friends: []bson.ObjectId{
			userCuteId,
		},
	}
	user.Update()
	monster.Update()

	room := model.Room{
		Type: model.ROOM_TYPE_FRIEND,
		Members: []bson.ObjectId{user.ID, monster.ID},
		CreateUser: user.ID,
	}

	room.Update()

	message := model.Message{
		UserID: user.ID,
		RoomID: room.ID,
		Type: model.MESSAGE_TYPE_CHAT,
		Content: "123123123",
	}

	message.Update()

	model.SetUserRoomMessageRead(room.ID.Hex(), user.ID.Hex())
}
