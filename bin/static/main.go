package main

import (
	"chat/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func main() {
	staticServer := gin.Default()

	staticController := StaticController{
		ServerPath: config.ENV_STATIC_SERVER_PATH,
		SavePath: config.ENV_STATIC_SAVE_PATH,
	}

	if _, err := os.Stat(staticController.SavePath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(staticController.SavePath, os.ModePerm); err != nil {
				panic(fmt.Sprintf("ERROR: mkdir static path failed, error: %s", err))
			}
		}
	}
	staticServer.Static(staticController.ServerPath, staticController.SavePath)
	staticServer.POST("/upload", staticController.Upload)

	staticServer.Run(config.ENV_STATIC_SERVER_URL)
}

type StaticController struct {
	ServerPath string
	SavePath string
}

func (controller *StaticController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	fileName := bson.NewObjectId().Hex()
	ext := filepath.Ext(file.Filename)
	if err := c.SaveUploadedFile(file, path.Join(controller.SavePath, fileName + ext)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"url": path.Join(controller.ServerPath, fileName + ext),
	})
}
