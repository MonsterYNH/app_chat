package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"time"
)

type Client struct {
	conn       *websocket.Conn
	updateTime int64
}

type WebSocketResponse struct {
	Ping     string `json:"ping"`
	Type     int    `json:"type"`
	Num      int    `json:"num"`
}

func CreateClient(conn *websocket.Conn) *Client {
	client := &Client{
		conn:       conn,
		updateTime: time.Now().Unix(),
	}
	return client
}

func (client *Client) ListenClient() {
	defer client.conn.Close()

	go client.heartBreak()
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			} else {
				log.Println("ERROR: read message failed, error: ", err)
				break
			}
		}
		if string(message) == "pong" {
			client.updateTime = time.Now().Unix()
		} else {
			log.Println("ERROR: wrong option with client")
			break
		}
	}
}

func (client *Client) heartBreak() {
	timer := time.NewTicker(time.Second)
	for {
		<-timer.C
		client.SendMessage(0,0, "ping")
		if time.Now().Add(-time.Second * 10).Unix() > client.updateTime {
			log.Println("客户端响应超时,断开")
			client.conn.Close()
			break
		}
	}
}

func (client *Client) SendMessage(kind, num int, option string) error {
	return client.conn.WriteJSON(WebSocketResponse{
		Ping: option,
		Type: kind,
		Num:num,
	})
}
