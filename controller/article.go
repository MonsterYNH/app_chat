package controller

import (
	"chat/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type ArticleController struct{}

type ArticleLikeEntry struct {
	Id     bson.ObjectId `json:"id"`
	Option string        `json:"option"`
}

func (controller *ArticleController) PostCreateArticle(c *gin.Context) {
	var article model.Article
	if err := c.BindJSON(&article); err != nil {
		model.Result(nil, 111, nil, err)
		return
	}
	userEntry, _ := c.Get("user")
	user := userEntry.(*model.User)
	article.UserId = user.ID
	if err := article.Update(); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, article, nil)
}

func (controller *ArticleController) GetArticleList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	userEntry, _ := c.Get("user")
	user := userEntry.(*model.User)

	articles, err := model.GetUserArticlesByPage(user.ID.Hex(), page, pageSize)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	for _, article := range articles {
		comments := make([]model.Comment, 0)
		for _, comment := range article.Comments {
			if comment.ToUserId != nil && comment.ToUserId.Valid() {
				userInfo, err := model.GetUserById(comment.ToUserId.Hex())
				if err != nil {
					model.Result(c, 111, nil, err)
					return
				}
				comment.ToUser = &model.UserInfo{
					Name:   userInfo.Name,
					Avatar: userInfo.Avatar,
				}
			}
			if comment.FromUserId != nil && comment.FromUserId.Valid() {
				userInfo, err := model.GetUserById(comment.FromUserId.Hex())
				if err != nil {
					model.Result(c, 111, nil, err)
					return
				}
				comment.FromUser = &model.UserInfo{
					Name:   userInfo.Name,
					Avatar: userInfo.Avatar,
				}
			}
			comments = append(comments, comment)
			article.Comments = comments
		}
	}
	bytes, _ := json.Marshal(articles)
	fmt.Println(string(bytes))
	model.Result(c, 111, articles, nil)
}

func (controller *ArticleController) PostCreateArticleComment(c *gin.Context) {
	var comment model.Comment
	if err := c.BindJSON(&comment); err != nil {
		model.Result(nil, 111, nil, err)
		return
	}
	if len(comment.Content) == 0 {
		model.Result(c, 111, nil, errors.New("ERROR: content can not be empty"))
		return
	}
	id := c.Param("id")
	article, err := model.GetArticleById(id)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	userEntry, _ := c.Get("user")
	fromUser := userEntry.(*model.User)
	comment.FromUserId = &fromUser.ID
	comment.ID = bson.NewObjectId()
	comment.CreateTime = time.Now().Unix()
	comment.ToUser = nil
	article.Comments = append(article.Comments, comment)
	if err := article.Update(); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	comment.FromUser = &model.UserInfo{
		Name: fromUser.Name,
		Avatar: fromUser.Avatar,
	}
	model.Result(c, 111, article, nil)
}

func (controller *ArticleController) LikeAndDislikeArticle(c *gin.Context) {
	var option ArticleLikeEntry
	if err := c.BindJSON(&option); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	if !option.Id.Valid() {
		model.Result(c, 111, nil, errors.New("ERROR: id is not valid"))
		return
	}
	article, err := model.GetArticleById(option.Id.Hex())
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	userEntry, _ := c.Get("user")
	user := userEntry.(*model.User)
	// 检查是否点赞
	isLike := false
	likeIndex := -1
	for index, entry := range article.LikeUserIds {
		if entry.Hex() == user.ID.Hex() {
			isLike = true
			likeIndex = index
			break
		}
	}
	ids := make([]bson.ObjectId, 0)
	switch (option.Option) {
	case "like":
		if !isLike {
			ids = append(article.LikeUserIds, user.ID)
		}
	case "dislike":
		if isLike {
			ids = append(article.LikeUserIds[:likeIndex], article.LikeUserIds[likeIndex+1:]...)
		}
	default:
		model.Result(c, 111, nil, errors.New("ERROR: wrong option"))
	}
	article.LikeUserIds = ids
	if err := article.Update(); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, gin.H{
		"like_user_ids": article.LikeUserIds,
	}, nil)
}
