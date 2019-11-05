package controller

import (
	"chat/model"
	"chat/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type UserController struct {}

func (controller *UserController) Login(c *gin.Context) {
	request := model.User{}
	if err := c.BindJSON(&request); err != nil {
		model.Result(nil, 111, nil, err)
		return
	}
	user, err := model.GetUserByAccountAndPassword(request.Account, request.Password)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	token, err := util.CreateToken(user.ID.Hex())
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, struct {
		*model.User
		Token string `json:"token"`
	}{
		User: user,
		Token: token,
	}, nil)
}

func (controller *UserController) Regist(c *gin.Context) {
	request := model.User{}
	if err := c.BindJSON(&request); err != nil {
		model.Result(nil, 111, nil, err)
		return
	}
	user := &model.User{
		Account: request.Account,
		Password: request.Password,
	}
	if err := user.Update(); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	token, err := util.CreateToken(user.ID.Hex())
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, struct {
		*model.User
		Token string `json:"token"`
	}{
		User: user,
		Token: token,
	}, nil)
}

func (controller *UserController) GetUserRelationShip(c *gin.Context) {
	userEntry, _ := c.Get("user")
	user := userEntry.(*model.User)
	users, err := user.GetUserFriends()
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, users, nil)
}

func (controller *UserController) UpdateUser(c *gin.Context) {
	userEntry, _ := c.Get("user")
	user := userEntry.(*model.User)
	update := model.User{}
	if err := c.BindJSON(&update); err != nil {
		model.Result(nil, 111, nil, err)
		return
	}
	if len(update.Password) != 0 && update.Password != user.Password {
		user.Password = update.Password
	}
	if len(update.Name) != 0 && update.Name != user.Name {
		user.Name = update.Name
	}
	if update.Sex != 0 && update.Sex != user.Sex {
		user.Sex = update.Sex
	}
	if update.Age != 0 && user.Age != update.Age {
		user.Age = update.Age
	}
	if len(update.Signture) != 0 && update.Signture != user.Signture {
		user.Signture = update.Signture
	}
	if len(update.Avatar) != 0 && update.Avatar != user.Avatar {
		user.Avatar = update.Avatar
	}
	if len(update.Friends) != 0 {
		for _, entry := range update.Friends {
			for _, entryFriend := range user.Friends {
				if entry.Hex() == entryFriend.Hex() {
					model.Result(c, 111, nil, errors.New(fmt.Sprintf("ERROR: %s is already your friend", entry.Hex())))
					return
				}
			}
			user.Friends = append(user.Friends, entry)
		}
	}
	if err := user.Update(); err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, user, nil)
}

func (controller *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := model.GetUserById(id)
	if err != nil {
		model.Result(c, 111, nil, err)
		return
	}
	model.Result(c, 111, model.UserInfo{
		Name: user.Name,
		Avatar: user.Avatar,
	}, nil)
}
