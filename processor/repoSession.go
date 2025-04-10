package processor

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
			return errors.WithStack(err)
		}
		util.ListReplacer(&content, "https://raw.githubusercontent.com", database.SettingMap["rand"][1])
	case "gitee":
		content, err = util.GetRepo("https://gitee.com/api/v5/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return errors.WithStack(err)
		}
	default:
		err = errors.WithStack(errors.New(info.RepoAPI + " is not supported"))
	}
	jsonString, err := json.Marshal(content)
	if err != nil {
		return err
	}
	repoDB := database.Repo{
		UUID:    uuid.New().String(),
		API:     info.RepoAPI,
		Info:    info.RepoInfo,
		Content: string(jsonString),
	}
	return errors.WithStack(repoDB.Create())
}

func refreshRepo(info *Session, data *database.Session) error {
	repoFinder := database.Repo{
		UUID: info.RepoUUID,
	}
	repoDBList, err := repoFinder.Read()
	if err != nil {
		return errors.WithStack(err)
	}
	repoDB := (*repoDBList).RepoList[0]
	info.RepoAPI = repoDB.API
	info.RepoInfo = repoDB.Info
	var content []string
	switch info.RepoAPI {
	case "github":
		content, err = util.GetRepo("https://api.github.com/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return errors.WithStack(err)
		}
		util.ListReplacer(&content, "https://raw.githubusercontent.com", database.SettingMap["rand"][1])
	case "gitee":
		content, err = util.GetRepo("https://gitee.com/api/v5/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return errors.WithStack(err)
		}
	default:
		err = errors.WithStack(errors.New(info.RepoAPI + " is not supported"))
	}
	jsonString, err := json.Marshal(content)
	if err != nil {
		return errors.WithStack(err)
	}
	repoDB.Content = string(jsonString)
	return errors.WithStack(repoDB.Update())
}

func delRepo(info *Session, data *database.Session) error {
	repoDB := database.Repo{
		UUID: info.RepoUUID,
	}
	return errors.WithStack(repoDB.Delete())
}

func fetchRepo(info *Session, data *database.Session) error {
	repoFinder := database.Repo{}
	repoDBList, err := repoFinder.Read()
	info.RepoData = repoDBList.RepoList
	if err != gorm.ErrRecordNotFound && err != nil {
		return errors.WithStack(err)
	}
	return nil
}
