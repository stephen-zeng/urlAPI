package processor

import (
	"errors"
	"github.com/google/uuid"
	"time"
	"urlAPI/database"
)

func login(info *Session, data *database.Session) error {
	var session database.Session
	if info.Operation == "login" && database.SettingMap["dash"][0] == data.Token {
		session.Token = uuid.New().String() + uuid.New().String()
		info.SessionToken = session.Token
		session.Term = info.LoginTerm
		if info.LoginTerm {
			session.Expire = time.Now().AddDate(0, 0, 7)
		} else {
			session.Expire = time.Now().AddDate(0, 0, 1)
		}
		if err := session.Create(); err != nil {
			return err
		}
		return nil
	}
	var ok bool
	session, ok = database.SessionMap[data.Token]
	switch {
	case !ok:
		return errors.New("Authentication failed")
	case time.Now().After(session.Expire):
		return errors.New("Expired token")
	case ok && time.Now().Before(session.Expire):
		return nil
	default:
		return errors.New("Authentication failed")
	}
	return nil
}

func logout(info *Session, data *database.Session) error {
	if err := data.Delete(); err != nil {
		return err
	}
	return nil
}

func exit(info *Session, data *database.Session) error {
	session, _ := database.SessionMap[data.Token]
	if !session.Term {
		if err := data.Delete(); err != nil {
			return err
		}
	}
	return nil
}
