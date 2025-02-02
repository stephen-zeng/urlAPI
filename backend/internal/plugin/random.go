package plugin

import (
	"backend/internal/data"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func random(API, User, Repo string) (PluginResponse, error) {
	var url string
	if API == "github" {
		url = "https://api.github.com/repos/" + User + "/" + Repo + "/contents"
	} else if API == "gitee" {
		url = "https://gitee.com/api/v5/repos/" + User + "/" + Repo + "/contents"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PluginResponse{}, err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return PluginResponse{}, err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return PluginResponse{}, errors.New("Failed to fetch repo contents")
	}
	var content []map[string]interface{}
	err = json.Unmarshal(jsonResponse, &content)
	if err != nil {
		return PluginResponse{}, err
	}
	length := len(content)
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(length) % length
	ret := content[index]["download_url"].(string)
	if API == "github" {
		replace, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"img"})))
		if err != nil {
			return PluginResponse{}, err
		}
		ret = strings.ReplaceAll(ret, "https://raw.githubusercontent.com", replace[0][len(replace[0])-1])
		return PluginResponse{
			URL: ret,
		}, nil
	} else {
		return PluginResponse{
			URL:     content[index]["download_url"].(string),
			Context: API + ", " + User + "/" + Repo,
		}, nil
	}
}
