package processor

import (
	"encoding/json"
	"errors"
	"urlAPI/database"
	"urlAPI/util"
)

func newRepo(info *Session, data *database.Session) error {
	var err error
	var content []string
	switch info.RepoAPI {
	case "github":
		content, err = util.GetRepo("https://api.github.com/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return err
		}
		util.ListReplacer(&content, "https://raw.githubusercontent.com", database.SettingMap["rand"][1])
	case "gitee":
		content, err = util.GetRepo("https://gitee.com/api/v5/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return err
		}
	default:
		err = errors.New(info.RepoAPI + " is not supported")
	}
	jsonString, err := json.Marshal(content)
	if err != nil {
		return err
	}
	repoDB := database.Repo{
		API:     info.RepoAPI,
		Info:    info.RepoInfo,
		Content: string(jsonString),
	}
	err = repoDB.Create()
	return err
}

func refreshRepo(info *Session, data *database.Session) error {
	repoFinder := database.Repo{
		UUID: info.RepoUUID,
	}
	repoDBList, err := repoFinder.Read()
	if err != nil {
		return err
	}
	repoDB := (*repoDBList).RepoList[0]
	info.RepoAPI = repoDB.API
	info.RepoInfo = repoDB.Info
	var content []string
	switch info.RepoAPI {
	case "github":
		content, err = util.GetRepo("https://api.github.com/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return err
		}
		util.ListReplacer(&content, "https://raw.githubusercontent.com", database.SettingMap["rand"][1])
	case "gitee":
		content, err = util.GetRepo("https://gitee.com/api/v5/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return err
		}
	default:
		err = errors.New(info.RepoAPI + " is not supported")
	}
	jsonString, err := json.Marshal(content)
	if err != nil {
		return err
	}
	repoDB.Content = string(jsonString)
	err = repoDB.Update()
	return err
}

func delRepo(info *Session, data *database.Session) error {
	repoDB := database.Repo{
		UUID: info.RepoUUID,
	}
	err := repoDB.Delete()
	if err != nil {
		return err
	}
	return nil
}

func fetchRepo(info *Session, data *database.Session) error {
	repoFinder := database.Repo{}
	repoDBList, err := repoFinder.Read()
	if err != nil {
		return err
	}
	info.RepoData = repoDBList.RepoList
	return nil
}
