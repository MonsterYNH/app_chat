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
	User        *UserInfo       `json:"user" bson:"user"`
	Content     string          `json:"content" bson:"content"`
	Images      []string        `json:"images" bson:"images"`
	CreateTime  int64           `json:"create_time" bson:"create_time"`
	LikeUserIds []bson.ObjectId `json:"like_user_ids" bson:"like_user_ids"`
	Comments    []Comment       `json:"comments" bson:"comments"`
	Status      int             `json:"-" bson:"status"`
}

type Comment struct {
	ID         bson.ObjectId  `json:"id" bson:"_id"`
	FromUserId *bson.ObjectId `json:"from_user_id" bson:"from_user_id"`
	FromUser   *UserInfo      `json:"from_user" bson:"from_user"`
	ToUserId   *bson.ObjectId `json:"to_user_id" bson:"to_user_id"`
	ToUser     *UserInfo      `json:"to_user" bson:"to_user"`
	Content    string         `json:"content" bson:"content"`
	CreateTime int64          `json:"create_time" bson:"create_time"`
	Status     int            `json:"-" bson:"status"`
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

func GetUserArticlesByPage(id string, page, pageSize int) ([]*Article, error) {
	idObject := bson.ObjectIdHex(id)
	if !idObject.Valid() {
		return nil, errors.New("ERROR: user id is not valid")
	}

	session := db.GetMgoSession()
	defer session.Close()

	articles := make([]*Article, 0)
	return articles, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ARTICLE).Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"status":  0,
				"user_id": idObject,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "user",
				"foreignField": "_id",
				"localField":   "user_id",
				"as":           "user",
			},
		},
		bson.M{
			"$unwind": "$user",
		},
		bson.M{
			"$sort": bson.M{
				"create_time": 1,
			},
		},
		bson.M{
			"$skip": (page - 1) * pageSize,
		},
		bson.M{
			"$limit": pageSize,
		},
	}).All(&articles)
	// articles, session.DB(config.ENV_DB_NAME).C(config.ENV_COLL_ARTICLE).Find(bson.M{"status": 0, "user_id": idObject}).Sort("create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&articles)
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
