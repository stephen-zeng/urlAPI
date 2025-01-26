package data

import (
	_ "github.com/mattn/go-sqlite3" // 下划线表示初始化这个包的内容以便使用
	"time"
)

func AddTask(data Config) (string, error) {
	id, err := addTask(time.Now(), data.IP, data.Type, data.Status, data.Target)
	if err != nil {
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
	if data.Status != "" {
		err = editTask(by, byData, data.Status, "")
	} else if data.Return != "" {
		err = editTask(by, byData, data.Return, "")
	}
	if err != nil {
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
		return nil, err
	} else {
		return ret, nil
	}
}
