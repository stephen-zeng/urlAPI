package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AlibabaImgInput struct {
	Prompt string `json:"prompt"`
}
type AlibabaImgParameters struct {
	Size string `json:"size"`
	N    int    `json:"n"`
}
type AlibabaImg struct {
	Model      string               `json:"model"`
	Input      AlibabaImgInput      `json:"input"`
	Parameters AlibabaImgParameters `json:"parameters"`
}

func alibabaTxt(prompt, contxt, model string) (string, error) {
	_, token, err := fetchConfig("alibaba")
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
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonPayload))
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

func alibabaImg(prompt, model, size string) (string, error) {
	_, token, err := fetchConfig("alibaba")
	if err != nil {
		return "", err
	}
	imgInput := AlibabaImgInput{
		Prompt: prompt,
	}
	imgParameter := AlibabaImgParameters{
		Size: size,
		N:    1,
	}
	imgPayload := AlibabaImg{
		Model:      model,
		Input:      imgInput,
		Parameters: imgParameter,
	}
	jsonPayload, err := json.Marshal(imgPayload)
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-DashScope-Async", "enable")
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
	}
	var response map[string]interface{}
	if err := json.Unmarshal(jsonResponse, &response); err != nil {
		return "", errors.New("plugin.response.error")
	}
	id := response["output"].(map[string]interface{})["task_id"].(string)
	timer := time.NewTimer(time.Second * 30)
	timeout := make(chan bool)
	go func() {
		<-timer.C
		fmt.Println("Timeout")
		timeout <- true
	}()
	for response["output"].(map[string]interface{})["status"] == "PENDING" {
		time.Sleep(1 * time.Second)
		response = fetchImgTask(id, token)
	}
	if response["output"].(map[string]interface{})["status"] == "FAILED" {
		return "", errors.New("plugin.response.error")
	} else if response["output"].(map[string]interface{})["status"] == "SUCCEEDED" {
		ret, err := json.Marshal(response)
		if err != nil {
			return "", errors.New("plugin.response.error")
		}
		return string(ret), nil
	}
	<-timeout
	return "", errors.New("plugin.response.timeout")
}

func fetchImgTask(id, token string) map[string]interface{} {
	req, _ := http.NewRequest("GET", "https://dashscope.aliyuncs.com/api/v1/tasks/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	jsonResponse, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(jsonResponse, &response); err != nil {
		return nil
	}
	return response
}
