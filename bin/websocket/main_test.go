package main

import (
	"fmt"
	"testing"
)

//func TestWebSocket(t *testing.T) {
//	url := "ws://localhost:3456"
//	dialer := websocket.Dialer{}
//	ws, _, err := dialer.Dial(url, http.Header{
//		"Auth": []string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzI0MDcyNjAsImp0aSI6IjVkYjZmNDRhNjgwMzBiN2QyMTkzYjY5MiJ9.0XtvPTQkCOPEl-CKQorCE_BTqJiFOMGL21b5qDDTEqY"},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for {
//		_, message, err := ws.ReadMessage()
//		if err != nil {
//			log.Println("ERROR: read message failed, error: ", err)
//			break
//		}
//
//		data := model.WebSocketResponse{}
//		if err := json.Unmarshal(message, &data); err != nil {
//			log.Println("ERROR: parse message failed, error: ", err)
//			break
//		}
//		bytes, _ := json.Marshal(data)
//		fmt.Println(string(bytes))
//		if data.Ping == "ping" {
//			if err := ws.WriteMessage(1, []byte("pong")); err != nil {
//				log.Println("ERROR: write message failed, error: ", err)
//				break
//			}
//		}
//	}
//
//	ws.Close()//关闭连接
//}


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
