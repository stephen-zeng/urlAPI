package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func deepseekTxt(prompt, contxt, model string) (string, error) {
	_, token, err := fetchConfig("deepseek")
	if err != nil {
		return "", err
	}
	userMessage := TxtMessage{
		Role:    "user",
		Content: prompt,
	}
	developerMessage := TxtMessage{
		Role:    "developer",
		Content: contxt,
	}
	txtPayload := Txt{
		Model:   model,
		Message: []TxtMessage{userMessage, developerMessage},
	}
	jsonPayload, err := json.Marshal(txtPayload)
	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("plugin.response.error")
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.New("plugin.response.error")
	} else {
		return string(jsonResponse), nil
	}
}
