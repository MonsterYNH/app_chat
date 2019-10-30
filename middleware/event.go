package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type MessageData struct {
	UserId string `json:"user_id"`
	Type   int    `json:"type"`
	Num    int    `json:"num"`
}

func requestWebSocketApi(url string, data interface{}) error {
	client := &http.Client{}

	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteData))
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
