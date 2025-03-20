package Database

import (
	"errors"
	"time"
)

type Session struct {
	Token  string    `json:"token" gorm:"primary_key"`
	Expire time.Time `json:"expire"`
	Term   bool      `json:"term"`
}

func sessionInit() error {
	err := db.AutoMigrate(&Session{})
	if err != nil {
		return errors.Join(errors.New("Session Init"), err)
	}
	return nil
}

func (session *Session) Create() (*Session, error) {
	err := db.Create(&session).Error
	if err != nil {
		return nil, errors.Join(errors.New("Session create"), err)
	}
	return session, nil

}

func (session *Session) Update() (*Session, error) {
	err := db.Save(session).Error
	if err != nil {
		return nil, errors.Join(errors.New("Session update"), err)
	}
	return session, nil
}

func (session *Session) Read() (*[]Session, error) {
	var sessions []Session
	err := db.Find(&sessions).Error
	if err != nil {
		return nil, errors.Join(errors.New("Session read"), err)
	}
	return &sessions, nil
}

func (session *Session) Delete() (*Session, error) {
	err := db.Delete(session).Error
	if err != nil {
		return nil, errors.Join(errors.New("Session delete"), err)
	}
	return session, nil
}
