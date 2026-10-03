package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"urlAPI/internal/model"
)

func (adapter *SQLiteAdapter) CreateRepo(repo *model.Repo) error {
	var content []string
	if err := json.Unmarshal([]byte(repo.Content), &content); err != nil {
		return errors.WithStack(err)
	}
	if err := adapter.db.Create(repo).Error; err != nil {
		return errors.WithStack(err)
	}
	Repos.Set(repo.API, repo.Info, content)
	return nil
}

func (adapter *SQLiteAdapter) UpdateRepo(repo *model.Repo) error {
	if err := adapter.db.Save(repo).Error; err != nil {
		return errors.WithStack(err)
	}
	var tmp []string
	if err := json.Unmarshal([]byte(repo.Content), &tmp); err != nil {
		return errors.WithStack(err)
	}
	Repos.Set(repo.API, repo.Info, tmp)
	return nil
}

func (adapter *SQLiteAdapter) ReadRepo(repo model.Repo) (*model.DBList, error) {
	var repos []model.Repo
	var err error
	switch {
	case repo.UUID != "":
		err = adapter.db.Where("uuid = ?", repo.UUID).Find(&repos).Error
	case repo.API != "":
		err = adapter.db.Where("api=? AND info=?", repo.API, repo.Info).Find(&repos).Error
	default:
		err = adapter.db.Find(&repos).Error
	}
	if len(repos) == 0 {
		err = gorm.ErrRecordNotFound
	}
	ret := model.DBList{
		RepoList: repos,
	}
	return &ret, errors.WithStack(err)
}

func (adapter *SQLiteAdapter) DeleteRepo(repo *model.Repo) error {
	// Callers usually only know the UUID; resolve the cache key from the row.
	if repo.UUID != "" && repo.API == "" && repo.Info == "" {
		var existing model.Repo
		err := adapter.db.Where("uuid = ?", repo.UUID).Limit(1).Find(&existing).Error
		if err != nil {
			return errors.WithStack(err)
		}
		repo.API, repo.Info = existing.API, existing.Info
	}
	if err := adapter.db.Delete(repo).Error; err != nil {
		return errors.WithStack(err)
	}
	Repos.Delete(repo.API, repo.Info)
	return nil
}

func CreateRepo(repo *model.Repo) error               { return localDB.CreateRepo(repo) }
func UpdateRepo(repo *model.Repo) error               { return localDB.UpdateRepo(repo) }
func ReadRepo(repo model.Repo) (*model.DBList, error) { return localDB.ReadRepo(repo) }
func DeleteRepo(repo *model.Repo) error               { return localDB.DeleteRepo(repo) }
