package data

import (
	"errors"
	_ "github.com/mattn/go-sqlite3" // 下划线表示初始化这个包的内容以便使用
	"log"
	"time"
)

func InitTask(data Config) error {
	err := db.AutoMigrate(&Task{})
	if err != nil {
		return err
	}
	if data.Type == "restore" {
		err := db.Where("1 = 1").Delete(&Task{})
		if err.Error != nil {
			return err.Error
		} else {
			log.Println("Initialized Task")
			return nil
		}
	}
	return nil
}

func NewTask(data Config) (string, error) {
	if data.TaskStatus == "" {
		data.TaskStatus = "waiting"
	}
	id, err := addTask(time.Now(),
		data.API,
		data.Type,
		data.TaskIP,
		data.TaskStatus,
		data.TaskTarget,
		data.TaskRegion,
		data.TaskSize,
		data.TaskModel,
		data.TaskReferer,
		data.TaskDevice)
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
	} else if data.TaskTarget != "" {
		err = delTask("target", data.TaskTarget)
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
	} else if data.TaskTarget != "" {
		by = "target"
		byData = data.TaskTarget
	}
	err = editTask(by, byData, data.TaskStatus, data.TaskReturn, data.TaskSize, data.API, data.TaskRegion)
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
	if data.Type != "" {
		ret, err = fetchTask(data.Type, data.By)
	} else if data.UUID != "" {
		ret, err = fetchTask("uuid", data.UUID)
	} else if data.TaskTarget != "" {
		ret, err = fetchTask("target", data.TaskTarget)
	} else if data.TaskIP != "" {
		ret, err = fetchTask("ip", data.TaskIP)
	} else if data.Type != "" {
		ret, err = fetchTask("type", data.TaskIP)
	} else if data.TaskRegion != "" {
		ret, err = fetchTask("region", data.TaskRegion)
	} else {
		err = errors.New("No valid filter")
	}
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		return ret, nil
	}
}
