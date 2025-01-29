// plugin还是同步的
package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type OpenaiImg struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
	N      int    `json:"n"`
}

func openaiTxt(prompt, contxt, model string) (string, error) {
	url, token, err := fetchConfig("openai")
	if err != nil {
		log.Println(err)
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
		Model:    model,
		Messages: []TxtMessage{userMessage, developerMessage},
	}
	jsonPayload, err := json.Marshal(txtPayload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.Join(err, errors.New(resp.Status))
	} else {
		return string(jsonResponse), nil
	}
}

func openaiImg(prompt, model, size string) (string, error) {
	url, token, err := fetchConfig("openai")
	if err != nil {
		return "", err
	}
	imgPayload := OpenaiImg{
		Model:  model,
		Prompt: prompt,
		Size:   size,
		N:      1,
	}
	jsonPayload, err := json.Marshal(imgPayload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.Join(err, errors.New(resp.Status))
	} else {
		return string(jsonResponse), nil
	}
}
