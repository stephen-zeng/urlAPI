package data

import (
	"bufio"
	"embed"
	"log"
	"strings"
)

//go:embed setting.init
var Init embed.FS

func InitSetting(data Config) (string, error) {
	var names []string
	var edits [][]string
	if data.Type == "" && db.Migrator().HasTable(&Setting{}) {
		return "", nil
	}
	db.AutoMigrate(&Setting{})
	f, _ := Init.Open("setting.init")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(names) == 0 {
			names = strings.Split(scanner.Text(), ";")
			continue
		}
		tmp := strings.Split(scanner.Text(), ";")
		edits = append(edits, tmp)
	}
	defer f.Close()
	err := editSetting(names, edits, data.Type == "update")
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
