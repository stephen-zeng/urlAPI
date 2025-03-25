package database

import (
	"errors"
	"github.com/google/uuid"
)

func repoInit() error {
	err := db.AutoMigrate(&Repo{})
	if err != nil {
		return errors.Join(errors.New("Repo Init"), err)
	}
	return nil
}
func (repo *Repo) Create() error {
	repo.UUID = uuid.New().String()
	err := db.Create(repo).Error
	if err != nil {
		return errors.Join(errors.New("Repo Create"), err)
	}
	return nil
}

func (repo *Repo) Update() error {
	err := db.Save(repo).Error
	if err != nil {
		return errors.Join(errors.New("Repo Update"), err)
	}
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
		err = errors.New("No valid Filter")
	}
	if err != nil {
		return nil, errors.Join(errors.New("Repo Read"), err)
	}
	return &DBList{
		RepoList: repos,
	}, nil
}

func (repo *Repo) Delete() error {
	err := db.Delete(repo).Error
	if err != nil {
		return errors.Join(errors.New("Repo Delete"), err)
	}
	return nil
}
