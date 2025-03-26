package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

func Txt(endpoint, token, model, context, prompt string) (string, error) {
	if endpoint == "" || token == "" || model == "" || context == "" || prompt == "" {
		return "", errors.Join(errors.New("Util TxtAPI insufficient info"))
	}
	userMessage := TxtMessage{
		Role:    "user",
		Content: prompt,
	}
	developerMessage := TxtMessage{
		Role:    "developer",
		Content: context,
	}
	txtPayload := TxtPayload{
		Model:    model,
		Messages: []TxtMessage{userMessage, developerMessage},
	}
	jsonPayload, err := json.Marshal(txtPayload)
	if err != nil {
		return "", errors.Join(errors.New("Util TxtAPI"), err)
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil {
		return "", errors.Join(errors.New("Util TxtAPI"), err)
	}
	defer resp.Body.Close()
	var txtResp TxtResp
	jsonResponse, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(jsonResponse, &txtResp)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.Join(errors.New("Util TxtAPI"), err, errors.New(resp.Status))
	} else {
		return txtResp.Choices[0].Message.Content, nil
	}
}

func AlibabaImg(token, prompt, model, size string) ([]byte, string, error) {
	imgInput := AlibabaImgInput{
		Prompt: prompt,
	}
	imgParameter := AlibabaImgParameters{
		Size: size,
		N:    1,
	}
	imgPayload := AlibabaImgPayload{
		Model:      model,
		Input:      imgInput,
		Parameters: imgParameter,
	}
	jsonPayload, _ := json.Marshal(imgPayload)
	req, _ := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-DashScope-Async", "enable")
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil {
		return nil, "", errors.Join(errors.New("Util AlibabaImgAPI"), err)
	}
	defer resp.Body.Close()
	var response AlibabaImgResp
	jsonResponse, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		return nil, "", errors.Join(errors.New("Util AlibabaImgAPI"), err)
	}
	id := response.Output.TaskID

	timer := time.NewTimer(time.Second * 30)
	timeout := make(chan bool)
	go func() {
		<-timer.C
		log.Println("Times up")
		timeout <- true
	}()

	for status := response.Output.TaskStatus; status == "PENDING" || status == "RUNNING"; status = response.Output.TaskStatus {
		time.Sleep(1 * time.Second)
		err = json.Unmarshal(alibabaFetchImgTask(id, token), &response)
		if err != nil {
			timer.Stop()
			return nil, "", errors.Join(errors.New("Util AlibabaImgAPI"), err)
		}
	}
	timer.Stop()

	if response.Output.TaskStatus != "SUCCEEDED" {
		log.Println(response)
		return nil, "", errors.Join(errors.New("Util AlibabaImgAPI"), err)
	}
	actualPrompt := response.Output.Results[0].ActualPrompt
	ret, err := Downloader(response.Output.Results[0].URL)
	if err != nil {
		return nil, "", errors.Join(errors.New("Util AlibabaImgAPI"), err)
	} else {
		return ret, actualPrompt, nil
	}
}

func alibabaFetchImgTask(id, token string) []byte {
	req, _ := http.NewRequest("GET", "https://dashscope.aliyuncs.com/api/v1/tasks/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	jsonResponse, _ := io.ReadAll(resp.Body)
	return jsonResponse
}

func OpenaiImg(endpoint, token, prompt, model, size string) ([]byte, error) {
	imgPayload := OpenaiImgPayload{
		Model:  model,
		Prompt: prompt,
		Size:   size,
		N:      1,
	}
	jsonPayload, _ := json.Marshal(imgPayload)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("Util OpenaiImgAPI"), err)
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.Join(errors.New("Util OpenaiImgAPI"), err)
	}
	var response OpenaiImgResp
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		return nil, errors.Join(errors.New("Util OpenaiImgAPI"), err)
	}
	ret, err := Downloader(response.Data[0].URL)
	if err != nil {
		return nil, errors.Join(errors.New("Util OpenaiImgAPI"), err)
	}
	return ret, nil
}
