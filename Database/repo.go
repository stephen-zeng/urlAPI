package Database

import (
	"errors"
	"github.com/google/uuid"
)

type Repo struct {
	UUID    string `json:"uuid" gorm:"primary_key"`
	API     string `json:"api"`
	Info    string `json:"info"`
	Content string `json:"content"`
}

func repoInit() error {
	err := db.AutoMigrate(&Repo{})
	if err != nil {
		return errors.Join(errors.New("Repo Init"), err)
	}
	return nil
}

func (repo *Repo) Create() (*Repo, error) {
	repo.UUID = uuid.New().String()
	err := db.Create(repo).Error
	if err != nil {
		return nil, errors.Join(errors.New("Repo Create"), err)
	}
	return repo, nil
}

func (repo *Repo) Update() (*Repo, error) {
	err := db.Save(repo).Error
	if err != nil {
		return nil, errors.Join(errors.New("Repo Update"), err)
	}
	return repo, nil
}

func (repo *Repo) Read() (*[]Repo, error) {
	var ret []Repo
	var err error
	switch {
	case repo.UUID != "":
		err = db.Where("uuid = ?", repo.UUID).Find(&ret).Error
	case repo.API != "":
		err = db.Where("api=? AND info=?", repo.API, repo.Info).Find(&ret).Error
	default:
		err = errors.New("No valid Filter")
	}
	if err != nil {
		return nil, errors.Join(errors.New("Repo Read"), err)
	}
	return &ret, nil
}

func (repo *Repo) Delete() (*Repo, error) {
	err := db.Delete(repo).Error
	if err != nil {
		return nil, errors.Join(errors.New("Repo Delete"), err)
	}
	return repo, nil
}
