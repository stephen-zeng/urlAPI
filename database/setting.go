package database

import (
	"encoding/json"
	"errors"
)

func settingInit() error {
	err := db.AutoMigrate(&Setting{})
	if err != nil {
		return errors.Join(errors.New("Setting Init"), err)
	}
	return nil
}

func (setting *Setting) Create() error {
	err := db.Create(&setting).Error
	if err != nil {
		return errors.Join(errors.New("Setting create"), err)
	}
	return nil

}

func (setting *Setting) Update() error {
	err := db.Save(setting).Error
	if err != nil {
		return errors.Join(errors.New("Setting update"), err)
	}
	var tmp []string
	err = json.Unmarshal([]byte(setting.Value), &tmp)
	if err != nil {
		return errors.Join(errors.New("Setting update"), err)
	}
	SettingMap[setting.Name] = tmp
	return nil
}

func (setting *Setting) Read() (*DBList, error) {
	var settings []Setting
	err := db.Find(&settings).Error
	if err != nil {
		return nil, errors.Join(errors.New("Setting read"), err)
	}
	return &DBList{
		SettingList: settings,
	}, nil
}

func (setting *Setting) Delete() error {
	err := db.Delete(setting).Error
	if err != nil {
		return errors.Join(errors.New("Setting delete"), err)
	}
	return nil
}
