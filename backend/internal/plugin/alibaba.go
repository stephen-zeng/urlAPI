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
		log.Println(err)
		return "", err
	}
	userMessage := TxtMessage{
		Role:    "user",
		Content: prompt,
	}
	developerMessage := TxtMessage{
		Role:    "system",
		Content: contxt,
	}
	txtPayload := Txt{
		Model:    model,
		Messages: []TxtMessage{userMessage, developerMessage},
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

func alibabaImg(prompt, model, size string) (string, error) {
	_, token, err := fetchConfig("alibaba")
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.Join(err, errors.New(resp.Status))
	}
	var response map[string]interface{}
	var originResponse []byte
	if err := json.Unmarshal(jsonResponse, &response); err != nil {
		return "", err
	}
	id := response["output"].(map[string]interface{})["task_id"].(string)
	timer := time.NewTimer(time.Second * 30)
	timeout := make(chan bool)
	go func() {
		<-timer.C
		log.Println("Timeout")
		timeout <- true
	}()
	for status := response["output"].(map[string]interface{})["task_status"].(string); status == "PENDING" || status == "RUNNING"; status = response["output"].(map[string]interface{})["task_status"].(string) {
		time.Sleep(1 * time.Second)
		originResponse = fetchImgTask(id, token)
		err := json.Unmarshal(originResponse, &response)
		if err != nil {
			return "", err
		}
	}
	if response["output"].(map[string]interface{})["task_status"] == "FAILED" {
		return "", err
	} else if response["output"].(map[string]interface{})["task_status"] == "SUCCEEDED" {
		return string(originResponse), nil
	}
	<-timeout
	return "", errors.New("Requirement Timeout")
}

func fetchImgTask(id, token string) []byte {
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
	return jsonResponse
}
