package op

import (
	"github.com/pkg/errors"
	"time"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func login(info *Session, data *model.Session) error {
	var session model.Session
	if info.Operation == "login" && database.SettingsStore.Get().Security.DashboardPasswordHash == data.Token {
		token, err := util.NewSessionToken()
		if err != nil {
			return err
		}
		session.Token = token
		info.SessionToken = session.Token
		session.Term = info.LoginTerm
		if info.LoginTerm {
			session.Expire = time.Now().AddDate(0, 0, 7)
		} else {
			session.Expire = time.Now().AddDate(0, 0, 1)
		}
		if err := db.CreateSession(&session); err != nil {
			return err
		}
		return nil
	}
	var ok bool
	session, ok = database.Sessions.Get(data.Token)
	switch {
	case !ok:
		return errors.WithStack(errors.New("Authentication failed"))
	case time.Now().After(session.Expire):
		return errors.New("Expired token")
	case ok && time.Now().Before(session.Expire):
		return nil
	default:
		return errors.WithStack(errors.New("Authentication failed"))
	}
}

func logout(data *model.Session) error {
	if err := db.DeleteSession(data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func exit(data *model.Session) error {
	session, _ := database.Sessions.Get(data.Token)
	if !session.Term {
		if err := db.DeleteSession(data); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
