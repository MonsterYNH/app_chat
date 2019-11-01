package middleware

import (
	"chat/model"
	socket "chat/model/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SendUserRoomMiddleware(c *gin.Context) {
	c.Next()

	userEntry, _ := c.Get("user")
	user := userEntry.(*model.User)
	fmt.Println("asdasdasdasdad")
	socket.WebSocketUtil.SendUserRoomEvent(user.ID.Hex())
}
