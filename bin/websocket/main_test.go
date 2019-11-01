package main

import (
	model "chat/model/websocket"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
)

func TestWebSocket(t *testing.T) {
	url := "ws://122.51.77.180:3456"
	dialer := websocket.Dialer{}
	ws, _, err := dialer.Dial(url, http.Header{
		"Auth": []string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzI0ODg1MDAsImp0aSI6IjVkYjhlYzI4NGMzZTNhMDAwNjRmMTQxNyJ9.Tx91oadKXxcIP2MVMW3ZRX3kqxkCA2P4NVPPsHvAUVo"},
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("ERROR: read message failed, error: ", err)
			break
		}

		data := model.WebSocketResponse{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("ERROR: parse message failed, error: ", err)
			break
		}
		bytes, _ := json.Marshal(data)
		fmt.Println(string(bytes))
		if data.Ping == "ping" {
			if err := ws.WriteMessage(1, []byte("pong")); err != nil {
				log.Println("ERROR: write message failed, error: ", err)
				break
			}
		}
	}

	ws.Close()//关闭连接
}


func TestChan(t *testing.T) {
	chan1 := make(chan int)
	chan2 := make(chan int)

	go func(chan1, chan2 chan int) {
		for {
			select {
			case data := <- chan1:
				fmt.Println("chan1: ", data)
			case data := <- chan2:
				fmt.Println("chan2: ", data)
			}
		}

	}(chan1, chan2)


	for i := 0; i < 10; i++ {
		sendChanData(chan1, i)
	}

	for i := 0; i < 10; i++ {
		sendChanData(chan2, i)
	}

}

func sendChanData(c chan int, data int) {
	c <- data
}
