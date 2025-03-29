package database

import (
	"errors"
)

func (session *Session) Create() error {
	err := db.Create(&session).Error
	if err != nil {
		errors.Join(errors.New("Session create"), err)
	}
	SessionMap[session.Token] = *session
	return nil

}

func (session *Session) Update() error {
	err := db.Save(session).Error
	if err != nil {
		errors.Join(errors.New("Session update"), err)
	}
	SessionMap[session.Token] = *session
	return nil
}

func (session *Session) Read() (*DBList, error) {
	var sessions []Session
	err := db.Where("token=?", session.Token).Find(&sessions).Error
	if len(sessions) == 0 {
		err = errors.New("No sessions found")
	}
	ret := DBList{
		SessionList: sessions,
	}
	if err != nil {
		return &ret, errors.Join(errors.New("Session read"), err)
	}
	return &ret, nil
}

func (session *Session) Delete() error {
	err := db.Delete(session).Error
	if err != nil {
		return errors.Join(errors.New("Session delete"), err)
	}
	delete(SessionMap, session.Token)
	return nil
}
