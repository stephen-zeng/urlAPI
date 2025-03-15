package data

import (
	"embed"
	"encoding/json"
	"errors"
	"io"
	"log"
)

//go:embed setting.json
var Init embed.FS

type SettingConfig struct {
	Names []string   `json:"names"`
	Edits [][]string `json:"edits"`
}

func InitSetting(data Config) (string, error) {
	if data.Type == "" && db.Migrator().HasTable(&Setting{}) {
		return "", nil
	}
	db.AutoMigrate(&Setting{})
	f, _ := Init.Open("setting.json")
	d, err := io.ReadAll(f)
	if err != nil {
		return "", errors.Join(errors.New("Error while reading init setting"), err)
	}
	var initSettings SettingConfig
	err = json.Unmarshal(d, &initSettings)
	if err != nil {
		return "", errors.Join(errors.New("Error while unmarshal init setting"), err)
	}
	err = editSetting(initSettings.Names, initSettings.Edits, data.Type == "update")
	defer f.Close()
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		log.Println("Initialized Setting")
		return "123456", nil
	}
}

func EditSetting(data Config) error {
	err := editSetting(data.SettingName, data.SettingEdit, false)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func FetchSetting(data Config) ([][]string, error) {
	ret, err := fetchSetting(data.SettingName)
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		return ret, nil
	}
}
