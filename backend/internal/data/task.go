package data

import (
	_ "github.com/mattn/go-sqlite3" // 下划线表示初始化这个包的内容以便使用
	"log"
	"time"
)

func NewTask(data Config) (string, error) {
	id, err := addTask(time.Now(), data.IP, data.Type, "waiting", data.Target)
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		return id, nil
	}
}

func DelTask(data Config) error {
	var err error
	if data.UUID != "" {
		err = delTask("uuid", data.UUID)
	} else if data.Target != "" {
		err = delTask("target", data.Target)
	}
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func EditTask(data Config) error {
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
	err = editTask(by, byData, data.Status, data.Return, data.Size, data.API)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func FetchTask(data Config) ([]Task, error) {
	var ret []Task
	var err error
	if data.UUID != "" {
		ret, err = fetchTask("uuid", data.UUID)
	} else if data.Target != "" {
		ret, err = fetchTask("target", data.Target)
	} else if data.IP != "" {
		ret, err = fetchTask("ip", data.IP)
	} else {
		ret, err = fetchTask("none", data.Status)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		return ret, nil
	}
}
