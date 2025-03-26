package database

import (
	"gorm.io/gorm"
	"time"
)

var (
	dbPath     = "assets/database.db"
	db         *gorm.DB
	SettingMap map[string][]string
	PromptMap  = map[string]int{
		"laugh":    0,
		"poem":     1,
		"sentence": 2,
	}
)

type Interface interface { //丢进来一个struct，然后就可以用这些方法了
	Create() error
	Update() error
	Read() (*DBList, error)
	Delete() error
}

type Repo struct {
	UUID    string `json:"uuid" gorm:"primary_key"`
	API     string `json:"api"`
	Info    string `json:"info"`
	Content string `json:"content"`
}

type Session struct {
	Token  string    `json:"token" gorm:"primary_key"`
	Expire time.Time `json:"expire"`
	Term   bool      `json:"term"`
}

type Task struct {
	// all
	UUID    string    `json:"uuid" gorm:"primary_key"`
	Time    time.Time `json:"time"`
	IP      string    `json:"ip"`
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Target  string    `json:"target"`
	Return  string    `json:"return"`
	Region  string    `json:"region"`
	Referer string    `json:"referer"`
	Device  string    `json:"device"`
	//txt, img, web, rand
	API string `json:"api"`

	// txt, img
	Model string `json:"model"`

	// img
	Size string `json:"size"`
}

type Setting struct {
	Name  string `json:"name" gorm:"primary_key"`
	Value string `json:"value"`
}

type DBList struct {
	RepoList    []Repo
	TaskList    []Task
	SessionList []Session
	SettingList []Setting
}
