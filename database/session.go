package database

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (session *Session) Create() error {
	SessionMap[session.Token] = *session
	return errors.WithStack(db.Create(&session).Error)

}

func (session *Session) Update() error {
	SessionMap[session.Token] = *session
	return errors.WithStack(db.Save(session).Error)
}

func (session *Session) Read() (*DBList, error) {
	var sessions []Session
	err := db.Where("token=?", session.Token).Find(&sessions).Error
	if len(sessions) == 0 {
		err = gorm.ErrRecordNotFound
	}
	ret := DBList{
		SessionList: sessions,
	}
	return &ret, errors.WithStack(err)
}

func (session *Session) Delete() error {
	delete(SessionMap, session.Token)
	return errors.WithStack(db.Delete(session).Error)
}
