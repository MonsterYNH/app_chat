package middleware

import (
	"chat/model"
	"chat/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Auth")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Code: 111,
			Message: err.Error(),
		})
		return
	}
	user, err := model.GetUserById(claims.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Code: 111,
			Message: err.Error(),
		})
		return
	}
	c.Set("user", user)
	fmt.Println(user.ID)
	c.Next()
}


