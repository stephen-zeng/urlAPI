package data

import (
	"errors"
	"fmt"
	"log"
)

type Setting struct {
	name  string   `gorm:"primaryKey"`
	value []string `gorm:"type:text"`
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

func addSetting(name []string, value [][]string) error {
	for index, nameItem := range name {
		valueItem := value[index]
		err := db.Create(Setting{
			name:  nameItem,
			value: valueItem,
		})
		if err != nil {
			return errors.New("setting.add.error")
		}
	}
	return nil
}

func editSetting(name []string, value [][]string) error {
	for index, nameItem := range name {
		valueItem := value[index]
		err := db.Model(&Setting{}).Where("name = ?", nameItem).Updates(Setting{
			name:  nameItem,
			value: valueItem,
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
		ret = append(ret, fetch[0].value)
	}
	return ret, nil
}
