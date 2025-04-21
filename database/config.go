package database

import (
	"gorm.io/gorm"
	"time"
)

var (
	dbPath     = "assets/database.db"
	db         *gorm.DB
	SettingMap = make(map[string][]string)
	PromptMap  = map[string]int{
		"laugh":    0,
		"poem":     1,
		"sentence": 2,
	}
	RepoMap    = make(map[string][]string)
	SessionMap = make(map[string]Session)
)

type Repo struct {
	UUID    string `json:"uuid" gorm:"primaryKey"`
	API     string `json:"api"`
	Info    string `json:"info"`
	Content string `json:"content"`
}

type Session struct {
	Token  string    `json:"token" gorm:"primaryKey"`
	Expire time.Time `json:"expire"`
	Term   bool      `json:"term"`
}

type Task struct {
	// all
	UUID     string    `json:"uuid" gorm:"primaryKey"`
	Time     time.Time `json:"time"`
	IP       string    `json:"ip"`
	Type     string    `json:"type"`
	Status   string    `json:"status"`
	Target   string    `json:"target"`
	Return   string    `json:"return"`
	Region   string    `json:"region"`
	Referer  string    `json:"referer"`
	Device   string    `json:"device"`
	MoreInfo string    `json:"more_info" gorm:"more_info"`
	//txt, img, web, rand
	API string `json:"api"`

	// txt, img
	Model string `json:"model"`
	Temp  string `json:"temp"`

	// img
	Size string `json:"size"`
}

type Setting struct {
	Name  string `json:"name" gorm:"primaryKey"`
	Value string `json:"value"`
}

type SettingInit struct {
	Names []string   `json:"names"`
	Edits [][]string `json:"edits"`
}

type DBList struct {
	RepoList    []Repo
	TaskList    []Task
	SessionList []Session
	SettingList []Setting
}
