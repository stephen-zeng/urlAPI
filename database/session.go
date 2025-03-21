package database

import (
	"errors"
	"time"
)

func sessionInit() error {
	err := db.AutoMigrate(&Session{})
	if err != nil {
		return errors.Join(errors.New("Session Init"), err)
	}
	return nil
}

func (session *Session) Create() error {
	err := db.Create(&session).Error
	if err != nil {
		errors.Join(errors.New("Session create"), err)
	}
	return nil

}

func (session *Session) Update() error {
	err := db.Save(session).Error
	if err != nil {
		errors.Join(errors.New("Session update"), err)
	}
	return nil
}

func (session *Session) Read() (*DBList, error) {
	var sessions []Session
	err := db.Find(&sessions).Error
	if err != nil {
		return nil, errors.Join(errors.New("Session read"), err)
	}
	return &DBList{
		SessionList: sessions,
	}, nil
}

func (session *Session) Delete() error {
	err := db.Delete(session).Error
	if err != nil {
		return errors.Join(errors.New("Session delete"), err)
	}
	return nil
}
