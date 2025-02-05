package data

import (
	"log"
)

type RepoResponse struct {
	UUID    string   `json:"uuid" gorm:"primaryKey"`
	API     string   `json:"api"`
	Info    string   `json:"info"`
	Content []string `json:"content"`
}

func InitRepo(data Config) error {
	if data.Type != "restore" && db.Migrator().HasTable(&Repo{}) {
		return nil
	} else {
		db.AutoMigrate(&Repo{})
		err := db.Where("1=1").Delete(&Repo{})
		if err.Error != nil {
			return err.Error
		} else {
			log.Println("Initialized Repo")
			return nil
		}
	}
}

func NewRepo(data Config) error {
	content, err := scanRepo(data.API, data.RepoInfo)
	if err != nil {
		return err
	}
	return addRepo(data.API, data.RepoInfo, content)
}

func RefreshRepo(data Config) error {
	list, err := fetchRepo("uuid", []string{data.UUID})
	if err != nil {
		return err
	}
	api := list[0].API
	info := list[0].Info
	content, err := scanRepo(api, info)
	if err != nil {
		return err
	}
	return editRepo(data.UUID, content)
}

func FetchRepo(data Config) ([]RepoResponse, error) {
	var ret []RepoResponse
	var err error
	if data.By == "api&info" {
		ret, err = fetchRepo("api&info", []string{data.API, data.RepoInfo})
	} else {
		ret, err = fetchRepo("none", []string{})
	}
	if err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}

func DelRepo(data Config) error {
	return delRepo(data.UUID)
}
