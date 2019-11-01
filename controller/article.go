package controller

import (
	"chat/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type ArticleController struct {}

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
	model.Result(c, 111, articles, nil)
}

func (controller *ArticleController) PostCreateArticleComment(c *gin.Context) {
	var comment model.Comment
	if err := c.BindJSON(&comment); err != nil {
		model.Result(nil, 111, nil, err)
		return
	}
	id := c.Param("id")
	article, err := model.GetArticleById(id)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	if !comment.ToUserId.Valid() {
		model.Result(c, 111, nil, errors.New("ERROR: to user id is not valid"))
		return
	}
	if _, err := model.GetUserById(comment.ToUserId.Hex()); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	userEntry, _ := c.Get("user")
	fromUser := userEntry.(*model.User)
	comment.FromUserId = fromUser.ID
	comment.ID = bson.NewObjectId()
	article.Comments = append(article.Comments, comment)
	if err := article.Update(); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, article, nil)
}
