package database

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"urlAPI/internal/model"
)

func (adapter *SQLiteAdapter) CreateSession(session *model.Session) error {
	if err := adapter.db.Create(session).Error; err != nil {
		return errors.WithStack(err)
	}
	Sessions.Set(*session)
	return nil
}

func (adapter *SQLiteAdapter) UpdateSession(session *model.Session) error {
	if err := adapter.db.Save(session).Error; err != nil {
		return errors.WithStack(err)
	}
	Sessions.Set(*session)
	return nil
}

func (adapter *SQLiteAdapter) ReadSession(session model.Session) (*model.DBList, error) {
	var sessions []model.Session
	err := adapter.db.Where("token=?", session.Token).Find(&sessions).Error
	if len(sessions) == 0 {
		err = gorm.ErrRecordNotFound
	}
	ret := model.DBList{
		SessionList: sessions,
	}
	return &ret, errors.WithStack(err)
}

func (adapter *SQLiteAdapter) DeleteSession(session *model.Session) error {
	// Evict from the cache even if the row delete fails so a revoked token
	// can no longer authenticate.
	Sessions.Delete(session.Token)
	return errors.WithStack(adapter.db.Delete(session).Error)
}

func CreateSession(session *model.Session) error               { return localDB.CreateSession(session) }
func UpdateSession(session *model.Session) error               { return localDB.UpdateSession(session) }
func ReadSession(session model.Session) (*model.DBList, error) { return localDB.ReadSession(session) }
func DeleteSession(session *model.Session) error               { return localDB.DeleteSession(session) }
