package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Setting struct {
	// 是可导出的类型才能使用GORM
	Name  string `gorm:"primaryKey"`
	Value string
}

func settingInit() error {
	var err error
	if !db.Migrator().HasTable(&Setting{}) {
		err = db.AutoMigrate(&Setting{})
	}
	if err != nil {
		fmt.Println("Error creating settings table")
		log.Fatal(err)
		return errors.New("setting.init.error")
	}
	fmt.Println("Initialized settings table")
	return nil
}

func editSetting(name []string, value [][]string) error {
	for index, nameItem := range name {
		valueItem := value[index]
		jsonData, _ := json.Marshal(valueItem)
		err := db.Save(Setting{
			Name:  nameItem,
			Value: string(jsonData),
		})
		if err.Error != nil {
			return errors.New("setting.edit.error")
		}
	}
	return nil
}

func fetchSetting(name []string) ([][]string, error) {
	var fetch []Setting
	var ret [][]string
	for _, nameItem := range name {
		err := db.Where("name=?", nameItem).Find(&fetch)
		if err.Error != nil {
			return nil, err.Error
		}
		var tmp []string
		_ = json.Unmarshal([]byte(fetch[0].Value), &tmp)
		ret = append(ret, tmp)
	}
	return ret, nil
}
