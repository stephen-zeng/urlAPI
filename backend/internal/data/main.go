package data

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3" // 下划线表示初始化这个包的内容以便使用
	"gorm.io/gorm"
	"time"
)

type Config struct {
	UUID   string
	Status string
	Return string
	Target string
	IP     string
	Type   string
}
type Option func(*Config)

func WithUUID(id string) Option {
	return func(config *Config) {
		config.UUID = id
	}
}
func WithStatus(status string) Option {
	return func(config *Config) {
		config.Status = status
	}
}
func WithReturn(ret string) Option {
	return func(config *Config) {
		config.Return = ret
	}
}
func WithTarget(target string) Option {
	return func(config *Config) {
		config.Target = target
	}
}
func WithIP(ip string) Option {
	return func(config *Config) {
		config.IP = ip
	}
}
func WithType(t string) Option {
	return func(config *Config) {
		config.Type = t
	}
}

func DataConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

var (
	dbPath string = "assets/database.db"
	db     *gorm.DB
	err    error
)

func init() {
	err := connect()
	if err == nil {
		dbInit()
	}
}
func Data() {
	//for i := 0; i < 10; i++ {
	//	add(time.Now(),
	//		"127.0.0.1",
	//		"txt.generate",
	//		"pending",
	//		strconv.Itoa(i))
	//}
	getItem, err := fetch("none", "")
	fmt.Println(getItem, err)
}

func Add(data Config) (string, error) {
	id, err := add(time.Now(), data.IP, data.Type, data.Status, data.Target)
	if err != nil {
		return "", err
	} else {
		return id, nil
	}
}

func Del(data Config) error {
	var err error
	if data.UUID != "" {
		err = del("uuid", data.UUID)
	} else if data.Target != "" {
		err = del("target", data.Target)
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}

func Edit(data Config) error {
	var by string
	var byData string
	var err error
	if data.UUID != "" {
		by = "uuid"
		byData = data.UUID
	} else if data.Target != "" {
		by = "target"
		byData = data.Target
	}
	if data.Status != "" {
		err = edit(by, byData, data.Status, "")
	} else if data.Return != "" {
		err = edit(by, byData, data.Return, "")
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}

func Fetch(data Config) ([]Task, error) {
	var ret []Task
	var err error
	if data.UUID != "" {
		ret, err = fetch("uuid", data.UUID)
	} else {
		ret, err = fetch("none", data.Target)
	}
	if err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}
