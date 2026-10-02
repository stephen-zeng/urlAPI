package database

import (
	"urlAPI/internal/model"
)

var (
	dbPath = "assets/database.db"
	// PromptMap is read-only after initialisation.
	PromptMap = map[string]int{
		"laugh":    0,
		"poem":     1,
		"sentence": 2,
	}
	// Sessions and Repos are populated from the database at startup and kept
	// in sync by the adapter methods; they are safe for concurrent use.
	Sessions = newSessionCache()
	Repos    = newRepoCache()
)

type Repo = model.Repo
type Session = model.Session
type Task = model.Task
type AppSetting = model.AppSetting
type Provider = model.Provider
type ServiceConfig = model.ServiceConfig
type Prompt = model.Prompt
type ConfigListItem = model.ConfigListItem
type DBList = model.DBList
