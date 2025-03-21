// plugin还是同步的
package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"urlAPI/internal/client"
)

type OpenaiImg struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
	N      int    `json:"n"`
}

func openaiTxt(prompt, contxt, model string) (PluginResponse, error) {
	url, token, err := fetchConfig("openai")
	if err != nil {
		log.Println(err)
		return PluginResponse{}, err
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
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		log.Println(err)
		return PluginResponse{}, err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	var ret txtResp
	err = json.Unmarshal(jsonResponse, &ret)
	if err != nil || resp.StatusCode != http.StatusOK {
		return PluginResponse{}, errors.Join(err, errors.New(resp.Status))
	} else {
		return PluginResponse{
			Response:     ret.Choices[0].Message.Content,
			InitPrompt:   prompt,
			ActualPrompt: prompt,
			Context:      contxt,
		}, nil
	}
}

func openaiImg(prompt, model, size string) (PluginResponse, error) {
	url, token, err := fetchConfig("openai")
	if err != nil {
		return PluginResponse{}, err
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
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return PluginResponse{}, err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return PluginResponse{}, errors.Join(err, errors.New(resp.Status))
	} else {
		var response openaiImgResp
		err = json.Unmarshal(jsonResponse, &response)
		if err != nil {
			return PluginResponse{}, err
		}
		return PluginResponse{
			URL:          response.Data[0].URL,
			ActualPrompt: prompt,
			InitPrompt:   prompt,
		}, nil
	}
}
