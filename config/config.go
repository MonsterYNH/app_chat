package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ENV_DB_NAME = "chat"
	ENV_COLL_ROOM = "room"
	ENV_COLL_USER = "user"
	ENV_COLL_MESSAGE = "message"
	ENV_COLL_ARTICLE = "article"

	ENV_MONGO_URL = "localhost:27017"
	ENV_SERVER_URL = "0.0.0.0:6000"
	ENV_STATIC_SERVER_URL = "0.0.0.0:7000"
	ENV_STATIC_SAVE_PATH string
	ENV_STATIC_SERVER_PATH = "/media"
	ENV_TOKEN_SECRET = "secret"

	ENV_SOCKET_API_URL = "http://localhost:3456/message"
)

func init() {
	if path := os.Getenv("ENV_STATIC_SAVE_PATH"); len(path) != 0 {
		ENV_STATIC_SAVE_PATH = path
	} else {
		defaultStaticPath, err := filepath.Abs("media")
		if err != nil {
			panic(fmt.Sprintf("ERROR: static path init error: %s", err))
		}
		ENV_STATIC_SAVE_PATH = defaultStaticPath
	}
	if secret := os.Getenv("ENV_TOKEN_SECRET"); len(secret) != 0 {
		ENV_TOKEN_SECRET = secret
	}
	if mongo := os.Getenv("ENV_MONGO_URL"); len(mongo) != 0 {
		ENV_MONGO_URL = mongo
	}
	if socketApi := os.Getenv("ENV_SOCKET_API_URL"); len(socketApi) > 0 {
		ENV_SOCKET_API_URL = socketApi
	}
}
