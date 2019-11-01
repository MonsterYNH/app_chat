package model

import (
	"chat/config"
	"chat/db"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Article struct {
	ID          bson.ObjectId   `json:"id" bson:"_id"`
	UserId      bson.ObjectId   `json:"user_id" bson:"user_id"`
	Content     string          `json:"content" bson:"content"`
	Images      []string        `json:"images" bson:"images"`
	CreateTime  int64           `json:"create_time" bson:"create_time"`
	LikeUserIds []bson.ObjectId `json:"-" bson:"like_user_ids"`
	LikeUsers   []User          `json:"like_users" bson:"-"`
	Comments    []Comment       `json:"comments" bson:"comments"`
	Status      int             `json:"-" bson:"status"`
}

type Comment struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	FromUserId bson.ObjectId `json:"from_user_id" bson:"from_user_id"`
	FromUser   *User         `json:"from_user" bson:"-"`
	ToUserId   bson.ObjectId `json:"to_user_id" bson:"to_user_id"`
	ToUser     *User         `json:"to_user" bson:"to_user"`
	Content    string        `json:"content" bson:"content"`
	CreateTime int64         `json:"create_time" bson:"create_time"`
	Status     int           `json:"-" bson:"status"`
}

func (article *Article) Update() error {
	if !article.ID.Valid() {
		article.ID = bson.NewObjectId()
	}
	if !article.UserId.Valid() {
		return errors.New("ERROR: user id can not empty")
	}
	if article.CreateTime == 0 {
		article.CreateTime = time.Now().Unix()
	}
	if len(article.Content) == 0 && len(article.Images) == 0 {
		return errors.New("ERROR: content and images can not empty")
	}

	session := db.GetMgoSession()
	defer session.Close()

	_, err := session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ARTICLE).UpsertId(article.ID, article)
	return err
}

func GetUserArticlesByPage(id string, page, pageSize int) ([]Article, error) {
	idObject := bson.ObjectIdHex(id)
	if !idObject.Valid() {
		return nil, errors.New("ERROR: user id is not valid")
	}

	session := db.GetMgoSession()
	defer session.Close()

	articles := make([]Article, 0)
	return articles, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ARTICLE).Find(bson.M{"status": 0, "user_id": idObject}).Sort("create_time").Skip((page-1)*pageSize).Limit(pageSize).All(&articles)
}

func GetArticleById(id string) (*Article, error) {
	idObject := bson.ObjectIdHex(id)
	if !idObject.Valid() {
		return nil, errors.New("ERROR: user id is not vaild")
	}

	session := db.GetMgoSession()
	defer session.Close()

	var article Article
	return &article, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ARTICLE).FindId(idObject).One(&article)
}
