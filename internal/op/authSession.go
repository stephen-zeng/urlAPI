package op

import (
	"github.com/pkg/errors"
	"log"
	"time"
	"urlAPI/internal/auth"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func login(info *Session, data *model.Session) error {
	var session model.Session
	if info.Operation == "login" && passwordLogin(data.Token) {
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

// passwordLogin reports whether credential (the client-side SHA-256 of the
// password) is the dashboard password. A token that is already a valid
// session is not treated as a password, which also spares the dashboard's
// session re-validation an Argon2id computation.
func passwordLogin(credential string) bool {
	if session, ok := database.Sessions.Get(credential); ok && time.Now().Before(session.Expire) {
		return false
	}
	stored := database.SettingsStore.Get().Security.DashboardPasswordHash
	ok, needsRehash := auth.VerifyPassword(stored, credential)
	if !ok {
		return false
	}
	if auth.IsDefaultCredential(credential) {
		log.Println("WARNING: the dashboard is using the default password; change it now")
	}
	if needsRehash {
		upgradePasswordHash(stored, credential)
	}
	return true
}

// upgradePasswordHash replaces a legacy or outdated stored hash after a
// successful login. Failures are logged; the login itself still succeeds.
func upgradePasswordHash(stored, credential string) {
	if auth.IsLegacyHash(stored) {
		// Legacy verification is case-insensitive; hash the canonical
		// lower-case form the dashboard sends.
		credential = stored
	}
	hash, err := auth.HashPassword(credential)
	if err != nil {
		log.Printf("Password hash upgrade failed: %v", err)
		return
	}
	err = database.UpdateAppSettings(func(settings *util.AppSettings) error {
		if settings.Security.DashboardPasswordHash != stored {
			return nil // changed concurrently
		}
		settings.Security.DashboardPasswordHash = hash
		return nil
	})
	if err != nil {
		log.Printf("Password hash upgrade failed: %v", err)
		return
	}
	log.Println("Upgraded the stored dashboard password hash to Argon2id")
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
