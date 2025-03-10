package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
	"urlAPI/internal/server"
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

func alibabaTxt(prompt, contxt, model string) (PluginResponse, error) {
	_, token, err := fetchConfig("alibaba")
	if err != nil {
		log.Println(err)
		return PluginResponse{}, err
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
	resp, err := server.GlobalHTTPClient.Do(req)
	if err != nil {
		log.Println(err)
		return PluginResponse{}, err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	ret := make(map[string]interface{})
	err = json.Unmarshal(jsonResponse, &ret)
	if err != nil || resp.StatusCode != http.StatusOK {
		return PluginResponse{}, errors.Join(err, errors.New(resp.Status))
	} else {
		return PluginResponse{
			Response:     ret["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string),
			InitPrompt:   prompt,
			ActualPrompt: prompt,
			Context:      contxt,
		}, nil
	}
}

func alibabaImg(prompt, model, size string) (PluginResponse, error) {
	_, token, err := fetchConfig("alibaba")
	if err != nil {
		log.Println(err)
		return PluginResponse{}, err
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
	resp, err := server.GlobalHTTPClient.Do(req)
	if err != nil {
		log.Println(err)
		return PluginResponse{}, err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return PluginResponse{}, errors.Join(err, errors.New(resp.Status))
	}
	response := make(map[string]interface{})
	if err := json.Unmarshal(jsonResponse, &response); err != nil {
		return PluginResponse{}, err
	}
	id := response["output"].(map[string]interface{})["task_id"].(string)
	timer := time.NewTimer(time.Second * 30)
	timeout := make(chan bool)
	go func() {
		<-timer.C
		log.Println("Times up")
		timeout <- true
	}()
	for status := response["output"].(map[string]interface{})["task_status"].(string); status == "PENDING" || status == "RUNNING"; status = response["output"].(map[string]interface{})["task_status"].(string) {
		time.Sleep(1 * time.Second)
		err := json.Unmarshal(fetchImgTask(id, token), &response)
		if err != nil {
			timer.Stop()
			return PluginResponse{}, err
		}
	}
	timer.Stop()
	if response["output"].(map[string]interface{})["task_status"] == "FAILED" {
		log.Println(response)
		return PluginResponse{}, errors.New("Alibaba imgGen Failed")
	} else if response["output"].(map[string]interface{})["task_status"] == "SUCCEEDED" {
		ret := response["output"].(map[string]interface{})["results"].([]interface{})[0].(map[string]interface{})
		return PluginResponse{
			URL:          ret["url"].(string),
			InitPrompt:   ret["orig_prompt"].(string),
			ActualPrompt: ret["actual_prompt"].(string),
		}, nil
	}
	<-timeout
	return PluginResponse{}, errors.New("Requirement Timeout")
}

func fetchImgTask(id, token string) []byte {
	req, _ := http.NewRequest("GET", "https://dashscope.aliyuncs.com/api/v1/tasks/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := server.GlobalHTTPClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	jsonResponse, _ := io.ReadAll(resp.Body)
	return jsonResponse
}
