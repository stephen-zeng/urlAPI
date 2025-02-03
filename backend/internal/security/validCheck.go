package security

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// 切片不能是常量
var txtAPI = []string{"openai", "alibaba", "deepseek", "otherapi"}
var imgAPI = []string{"openai", "alibaba"}

func modelCheck(Type, API string) error {
	if Type == "txt.gen" {
		for _, api := range txtAPI {
			if API == api {
				return nil
			}
		}
		log.Println("The request " + API + " is NOT valid.")
		return errors.New("invalid request")
	} else if Type == "img.gen" {
		for _, api := range imgAPI {
			if API == api {
				return nil
			}
		}
		log.Println("The request " + API + " is NOT valid.")
		return errors.New("invalid request")
	} else {
		log.Println("The request " + API + " is NOT valid.")
		return errors.New("invalid request")
	}
}

func repoCheck(API, User, Repo string) error {
	var url string
	if API == "github" {
		url = "https://api.github.com/repos/" + User + "/" + Repo
	} else if API == "gitee" {
		url = "https://gitee.com/api/v5/repos/" + User + "/" + Repo
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(User + "/" + Repo + " is not valid.")
		return errors.New("repoCheck failed")
	} else {
		return nil
	}
}
