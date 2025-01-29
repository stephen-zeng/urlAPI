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

func random(API, User, Repo string) (string, error) {
	var url string
	if API == "github" {
		url = "https://api.github.com/repos/" + User + "/" + Repo + "/contents"
	} else if API == "gitee" {
		url = "https://gitee.com/api/v5/repos/" + User + "/" + Repo + "/contents"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
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
		return "", errors.New("Failed to fetch repo contents")
	}
	var content []map[string]interface{}
	err = json.Unmarshal(jsonResponse, &content)
	if err != nil {
		return "", err
	}
	length := len(content)
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(length) % length
	ret := content[index]["download_url"].(string)
	if API == "github" {
		replace, err := data.FetchSetting(data.DataConfig(data.WithName("rand")))
		if err != nil {
			return "", err
		}
		ret = strings.ReplaceAll(ret, "https://raw.githubusercontent.com", replace[0][0])
		return ret, nil
	}
	return content[index]["download_url"].(string), nil
}
