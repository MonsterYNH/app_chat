package db

import (
	"chat/config"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
)

var (
	MongoSession *mgo.Session
)

func init() {
	log.Println("MongoDB url:", config.ENV_MONGO_URL)
	var err error
	MongoSession, err = mgo.Dial(config.ENV_MONGO_URL)
	if err != nil {
		panic(fmt.Sprintf("Error: connect to mongo db failed, error: %s", err))
	}

	if err := MongoSession.Ping(); err != nil {
		panic(fmt.Sprintf("Error: ping to mongo db failed, error: %s", err))
	}
	log.Println("MongoDB init success")
}

func GetMgoSession() *mgo.Session {
	return MongoSession.Copy()
}

func FlushDb() {
	if err := MongoSession.DB(config.ENV_DB_NAME).DropDatabase(); err != nil {
		log.Println("ERROR: flush db failed, error: ", err)
	} else {
		log.Println("ERROR: flush db success")
	}
}
