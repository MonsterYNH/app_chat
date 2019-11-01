package model

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Result(c *gin.Context, code int, data interface{}, err error) {
	message := "Success"
	if err != nil {
		message = err.Error()
		log.Println("ERROR: ", err)
	}
	if c != nil {
		c.JSON(http.StatusOK, Response{
			Code: code,
			Message: message,
			Data: data,
		})
	}
}
