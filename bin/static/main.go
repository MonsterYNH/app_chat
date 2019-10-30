package main

import (
	"chat/config"
	"chat/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
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
		model.Result(c, 111, nil, err)
		return
	}
	fileName := bson.NewObjectId().Hex()
	ext := filepath.Ext(file.Filename)
	if err := c.SaveUploadedFile(file, path.Join(controller.SavePath, fileName + ext)); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, gin.H{
		"url": path.Join(controller.ServerPath, fileName + ext),
	}, nil)
}
