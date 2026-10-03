package op

import (
	"encoding/json"
	"net/url"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

// fetchRepoContent lists the files of a GitHub or Gitee repository given as
// "owner/repo" and returns the filtered download links.
func fetchRepoContent(api, info string) ([]string, error) {
	owner, repo, err := util.SplitRepoInfo(info)
	if err != nil {
		return nil, err
	}
	path := url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/contents"
	var content []string
	switch api {
	case "github":
		content, err = util.GetRepo("https://api.github.com/repos/" + path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		util.ListReplacer(&content, "https://raw.githubusercontent.com", database.SettingsStore.Get().Random.SourceRewriteFrom)
	case "gitee":
		content, err = util.GetRepo("https://gitee.com/api/v5/repos/" + path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	default:
		return nil, errors.WithStack(errors.New(api + " is not supported"))
	}
	return util.LinkFilter(content), nil
}

func newRepo(info *Session) error {
	content, err := fetchRepoContent(info.RepoAPI, info.RepoInfo)
	if err != nil {
		return err
	}
	jsonString, err := json.Marshal(content)
	if err != nil {
		return err
	}
	repoDB := model.Repo{
		UUID:    uuid.New().String(),
		API:     info.RepoAPI,
		Info:    info.RepoInfo,
		Content: string(jsonString),
	}
	return errors.WithStack(db.CreateRepo(&repoDB))
}

func refreshRepo(info *Session) error {
	repoFinder := model.Repo{
		UUID: info.RepoUUID,
	}
	repoDBList, err := db.ReadRepo(repoFinder)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(repoDBList.RepoList) == 0 {
		return errors.New("repository not found")
	}
	repoDB := repoDBList.RepoList[0]
	info.RepoAPI = repoDB.API
	info.RepoInfo = repoDB.Info
	content, err := fetchRepoContent(info.RepoAPI, info.RepoInfo)
	if err != nil {
		return err
	}
	jsonString, err := json.Marshal(content)
	if err != nil {
		return errors.WithStack(err)
	}
	repoDB.Content = string(jsonString)
	return errors.WithStack(db.UpdateRepo(&repoDB))
}

func delRepo(info *Session) error {
	repoDB := model.Repo{
		UUID: info.RepoUUID,
	}
	return errors.WithStack(db.DeleteRepo(&repoDB))
}

func fetchRepo(info *Session) error {
	repoFinder := model.Repo{}
	repoDBList, err := db.ReadRepo(repoFinder)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return errors.WithStack(err)
	}
	if repoDBList != nil {
		info.RepoData = repoDBList.RepoList
	}
	return nil
}
