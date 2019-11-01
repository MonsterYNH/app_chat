package socket

import (
	"bytes"
	"chat/config"
	"chat/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	EVENT_USER_LOGIN_STASTUS = iota + 1
	EVENT_USER_ROOM_STATUS
	EVENT_USER_MESSAGE_STATUS
)

type SocketUtil struct {
	Url string
}

var WebSocketUtil = &SocketUtil{
	Url: config.ENV_SOCKET_API_URL,
}

func (util *SocketUtil) SendUserLoginEvent(userId string) {
	if err := util.baseUrl(userId, EVENT_USER_LOGIN_STASTUS, 0); err != nil {
		log.Println("ERROR: get user room unread message failed, error: ", err)
	}
}

func (util *SocketUtil) SendUserRoomEvent(userId string) {
	rooms, err := model.GetRooms(userId)
	if err != nil {
		log.Println("ERROR: get user rooms failed, error: ", err)
		return
	}

	for _, entry := range rooms {
		for _, user := range entry.Members {
			fmt.Println(user.Hex(), "========", userId)
			if user.Hex() != userId {
				unRead, err := model.GetUserRoomUnReadMessageCount(entry.ID.Hex(), user.Hex())
				if err != nil {
					log.Println("ERROR: get user room unread message failed, error: ", err)
					continue
				}
				if err := util.baseUrl(user.Hex(), EVENT_USER_ROOM_STATUS, unRead); err != nil {
					log.Println("ERROR: send user room unread message failed, error: ", err)
					continue
				}
			}
		}
	}
}

func (util *SocketUtil) baseUrl(userId string, kind, num int) error {
	client := &http.Client{}

	byteData, err := json.Marshal(map[string]interface{}{
		"user_id": userId,
		"type": kind,
		"num": num,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", util.Url, bytes.NewBuffer(byteData))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("ERROR: request with error: " + string(byteData))
	}
	return nil
}

