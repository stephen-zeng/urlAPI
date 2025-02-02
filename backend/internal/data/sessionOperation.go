package data

import (
	"errors"
	"time"
)

type Session struct {
	Token  string `gorm:"primary_key"`
	Expire time.Time
	Term   bool
}

func addSession(Token string, Term bool, Expire time.Time) error {
	err := db.Create(Session{
		Token:  Token,
		Expire: Expire,
		Term:   Term,
	})
	if err != nil {
		return err.Error
	} else {
		return nil
	}
}

func delSession(Token string) error {
	err := db.Where("token = ?", Token).Delete(&Session{})
	if err != nil {
		return err.Error
	} else {
		return nil
	}
}

func fetchSession(Token string) ([]Session, error) {
	var ret []Session
	err := db.Where("token = ?", Token).Find(&ret)
	if err != nil {
		return nil, err.Error
	}
	if len(ret) == 0 {
		return nil, errors.New("session not found")
	} else {
		return ret, nil
	}
}
