package utils

import (
	"InnerG/config"
	"InnerG/pkg/constants"
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

type ApiReq struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func SendMessageToAPI(ms []Message) (*http.Response, error) {
	if config.Api == nil {
		return nil, fmt.Errorf("SendMessageToAPI : API Config Empty")
	}
	reqBody := ApiReq{
		Model:    config.Api.Model,
		Stream:   true,
		Messages: ms,
	}
	bodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("SendMessageToAPI: parse json error: %w", err)
	}
	req, err := http.NewRequest(
		constants.ApiRequestWay,
		config.Api.Url,
		bytes.NewReader(bodyJson))

	if err != nil {
		return nil, fmt.Errorf("SendMessageToAPI: new http request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ProcessApiKey())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SendMessageToAPI: send http request error: %w", err)
	}
	return resp, nil
}

func ProcessApiKey() string {
	return "Bearer " + config.Api.Key
}
