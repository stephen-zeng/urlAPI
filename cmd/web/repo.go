package web

import (
	"encoding/json"
	"errors"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"urlAPI/internal/client"
	"urlAPI/internal/file"
)

func getRepoCount(x float64) string {
	if x >= 1000 {
		return strconv.FormatFloat(x/1000.0, 'f', 1, 64) + "k"
	} else {
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
}

func repo(URL string, From, UUID, Token string) (WebResponse, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if Token != "" {
		req.Header.Set("Authorization", "Bearer "+Token)
	}
	if err != nil {
		return WebResponse{}, err
	}
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return WebResponse{}, err
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting github repo info")
		return WebResponse{}, errors.Join(err, errors.New(resp.Status))
	}
	var repo repoResp
	err = json.Unmarshal(jsonResp, &repo)
	if err != nil {
		return WebResponse{}, err
	}
	author := repo.Owner.Login
	name := repo.Name
	description := repo.Description
	forkCount := getRepoCount(repo.ForksCount)
	starCount := getRepoCount(repo.StargazersCount)
	bgFile, err := file.LogoFS.Open("github_logo.png")
	bgImg, err := png.Decode(bgFile)
	if err != nil {
		log.Println("Unable to decode github background image")
		return WebResponse{}, err
	}
	err = drawRepo(bgImg, name, author, description, starCount, forkCount, UUID)
	if err != nil {
		log.Println("Error when drawing the img")
		return WebResponse{}, err
	}
	return WebResponse{
		Target: URL,
		URL:    From + "/download?img=" + UUID,
	}, err
}
