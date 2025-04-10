package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (repo *Repo) Create() error {
	return errors.WithStack(db.Create(repo).Error)
}

func (repo *Repo) Update() error {
	if err := db.Save(repo).Error; err != nil {
		return errors.WithStack(err)
	}
	var tmp []string
	if err := json.Unmarshal([]byte(repo.Content), &tmp); err != nil {
		return errors.WithStack(err)
	}
	RepoMap[repo.API+";"+repo.Info] = tmp
	return nil
}

func (repo *Repo) Read() (*DBList, error) {
	var repos []Repo
	var err error
	switch {
	case repo.UUID != "":
		err = db.Where("uuid = ?", repo.UUID).Find(&repos).Error
	case repo.API != "":
		err = db.Where("api=? AND info=?", repo.API, repo.Info).Find(&repos).Error
	default:
		err = db.Find(&repos).Error
	}
	if len(repos) == 0 {
		err = gorm.ErrRecordNotFound
	}
	ret := DBList{
		RepoList: repos,
	}
	return &ret, errors.WithStack(err)
}

func (repo *Repo) Delete() error {
	delete(RepoMap, repo.API+";"+repo.Info)
	return errors.WithStack(db.Delete(repo).Error)
}
