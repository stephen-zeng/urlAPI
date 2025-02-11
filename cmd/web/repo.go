package web

import (
	"embed"
	"encoding/json"
	"errors"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

//go:embed github_logo.png gitee_logo.png
var bgFS embed.FS

func repo(URL string, From, UUID, Token string) (WebResponse, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if Token != "" {
		req.Header.Set("Authorization", "Bearer "+Token)
	}
	if err != nil {
		return WebResponse{}, err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return WebResponse{}, err
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting github repo info")
		return WebResponse{}, errors.Join(err, errors.New(resp.Status))
	}
	var repo map[string]interface{}
	err = json.Unmarshal(jsonResp, &repo)
	if err != nil {
		return WebResponse{}, err
	}
	author := repo["owner"].(map[string]interface{})["login"].(string)
	name := repo["name"].(string)
	description := repo["description"].(string)
	forkCount := repo["forks_count"].(float64)
	starCount := repo["stargazers_count"].(float64)
	var fork, star string
	if forkCount >= 1000 {
		fork = strconv.FormatFloat(forkCount/1000.0, 'f', 1, 64) + "k"
	} else {
		fork = strconv.FormatFloat(forkCount, 'f', -1, 64)
	}
	if starCount >= 1000 {
		star = strconv.FormatFloat(starCount/1000.0, 'f', 1, 64) + "k"
	} else {
		star = strconv.FormatFloat(starCount, 'f', -1, 64)
	}
	bgFile, err := bgFS.Open("github_logo.png")
	bgImg, err := png.Decode(bgFile)
	if err != nil {
		log.Println("Unable to decode github background image")
		return WebResponse{}, err
	}
	err = drawRepo(bgImg, name, author, description, star, fork, UUID)
	if err != nil {
		log.Println("Error when drawing the img")
		return WebResponse{}, err
	}
	return WebResponse{
		URL: From + "/download?img=" + UUID,
	}, err
}
