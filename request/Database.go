package request

import (
	"urlAPI/database"
)

type DB struct {
	Task      database.Task
	Repo      database.Repo
	Setting   database.Setting
	Session   database.Session
	Operation database.Interface
}
