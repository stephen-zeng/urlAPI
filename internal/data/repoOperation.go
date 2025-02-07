package data

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Repo struct {
	UUID    string `json:"uuid" gorm:"primaryKey"`
	API     string `json:"api"`
	Info    string `json:"info"`
	Content string `json:"content"`
}

func addRepo(API, Info string, Content []string) error {
	id := uuid.New().String()
	jsonDat, _ := json.Marshal(Content)
	err := db.Create(Repo{
		UUID:    id,
		API:     API,
		Info:    Info,
		Content: string(jsonDat),
	})
	if err.Error != nil {
		log.Print(err.Error)
		return err.Error
	} else {
		return nil
	}
}

func editRepo(UUID string, Content []string) error {
	jsonDat, _ := json.Marshal(Content)
	err := db.Model(&Repo{}).Where("uuid = ?", UUID).Updates(Repo{
		Content: string(jsonDat),
	})
	if err.Error != nil {
		log.Print(err.Error)
		return err.Error
	} else {
		return nil
	}
}
func delRepo(UUID string) error {
	err := db.Where("uuid = ?", UUID).Delete(&Repo{})
	if err.Error != nil {
		log.Print(err.Error)
		return err.Error
	} else {
		return nil
	}
}
func fetchRepo(by string, data []string) ([]RepoResponse, error) {
	var repos []Repo
	var err *gorm.DB
	if by == "uuid" {
		err = db.Where("uuid = ?", data[0]).Find(&repos)
	} else if by == "api&info" {
		err = db.Where("api = ? AND info = ?", data[0], data[1]).Find(&repos)
	} else if by == "none" {
		err = db.Find(&repos)
	} else {
		return nil, errors.New("Unknown Filter")
	}
	if err.Error != nil {
		log.Print(err.Error)
		return nil, err.Error
	}
	if len(repos) == 0 {
		return nil, errors.New("Repo not found")
	}
	var ret []RepoResponse
	for _, repo := range repos {
		var tmp []string
		_ = json.Unmarshal([]byte(repo.Content), &tmp)
		ret = append(ret, RepoResponse{
			UUID:    repo.UUID,
			API:     repo.API,
			Info:    repo.Info,
			Content: tmp,
		})
	}
	return ret, nil
}

func scanRepo(API, Info string) ([]string, error) {
	var url string
	if API == "github" {
		url = "https://api.github.com/repos/" + Info + "/contents"
	} else if API == "gitee" {
		url = "https://gitee.com/api/v5/repos/" + Info + "/contents"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch repo contents")
	}
	var content []map[string]interface{}
	var ret []string
	var replace string
	err = json.Unmarshal(jsonResponse, &content)
	if err != nil {
		return nil, err
	}
	if API == "github" {
		list, err := FetchSetting(DataConfig(WithSettingName([]string{"rand"})))
		if err != nil {
			return nil, err
		}
		replace = list[0][1]
	}
	for _, item := range content {
		url := item["download_url"].(string)
		if API == "github" {
			url = strings.ReplaceAll(url, "https://raw.githubusercontent.com", replace)
		}
		ret = append(ret, url)
	}
	return ret, nil
}
