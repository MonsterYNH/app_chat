package model

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestRoom_Update(t *testing.T) {
	room := Room{}
	if err := room.Update(); err != nil {
		t.Fatal(err)
	}
	t.Log(room)
}

func TestGetRooms(t *testing.T) {
	rooms, err := GetRooms("5db035910d833f013a4167e1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rooms)
}

func TestGetRoomUsers(t *testing.T) {
	users, err := GetRoomUsers("5db035910d833f013a4167e1")
	if err != nil {
		t.Fatal(err)
	}
	bytes, _ := json.Marshal(users)
	t.Log(string(bytes))
}

func TestUser_Update(t *testing.T) {
	for i := 0; i < 10; i++ {
		user := User{
			Account: "monster1",
			Password: "monster1",
		}
		if err := user.Update(); err != nil {
			t.Fatal(err)
		}

		update, err := GetUserById("5db1b32368030bb3c50ba869")
		if err != nil {
			t.Fatal(err)
		}
		update.Friends = append(update.Friends, user.ID)
		if err := update.Update(); err != nil {
			t.Fatal(err)
		}
		t.Log(user)
	}

}

func TestGetUserById(t *testing.T) {
	user, err := GetUserById("5db111d60d833f0ce47f8559")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestMessage_Update(t *testing.T) {
	message := Message{
		RoomID: bson.NewObjectId(),
		UserID: bson.NewObjectId(),
		Content: "123",
	}
	if err := message.Update(); err != nil {
		t.Fatal(err)
	}
	t.Log(message)
}

func TestArticle_Update(t *testing.T) {
	article := Article{
		UserId: bson.NewObjectId(),
		Content: "第一篇文章",
	}
	if err := article.Update(); err != nil {
		t.Fatal(err)
	}
}

func TestGetArticleById(t *testing.T) {
	ids := []int{1,2,3,4,5,6,7,8,9}
	index := 2
	ids = append(ids[:index], ids[index+1:]...)
	fmt.Println(ids)
}


