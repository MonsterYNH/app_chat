package main

import (
	"chat/model"
	socket "chat/model/websocket"
	"chat/util"
	"encoding/json"
	"fmt"

	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	engine := gin.Default()
	engine.GET("/", WebSocket)
	message := MessageController{}
	engine.POST("/message", message.SendMessage)
	engine.Run("0.0.0.0:3456")
}

type MessageController struct{}

var Manage = socket.CreateManage()

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Data struct {
	UserId string `json:"user_id"`
	Type   int    `json:"type"`
	Num    int    `json:"num"`
}

func WebSocket(c *gin.Context) {
	token := c.GetHeader("Auth")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Code:    111,
			Message: err.Error(),
		})
		return
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ERRPR: upgrader websocket failed, error:", err)
		return
	}
	defer ws.Close()

	client := socket.CreateClient(ws)
	Manage.AddClient(claims.Id, client)
	client.ListenClient()
}

func (controller *MessageController) SendMessage(c *gin.Context) {
	var messageData Data
	if err := c.BindJSON(&messageData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	bytes, _ := json.Marshal(messageData)
	fmt.Println(string(bytes))
	if len(messageData.UserId) == 0 || messageData.Type == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errors.New("ERROR: parameter error"),
		})
		return
	}
	client, exist := Manage.GetClient(messageData.UserId)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "client is out of line",
		})
		return
	}
	if err := client.SendMessage(messageData.Type, messageData.Num, "ping"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
