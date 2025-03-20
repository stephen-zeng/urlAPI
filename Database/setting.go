package Database

import (
	"errors"
)

type Setting struct {
	Name  string `json:"name" gorm:"primary_key"`
	Value string `json:"value"`
}

func settingInit() error {
	err := db.AutoMigrate(&Setting{})
	if err != nil {
		return errors.Join(errors.New("Setting Init"), err)
	}
	return nil
}

func (setting *Setting) Create() (*Setting, error) {
	err := db.Create(&setting).Error
	if err != nil {
		return nil, errors.Join(errors.New("Setting create"), err)
	}
	return setting, nil

}

func (setting *Setting) Update() (*Setting, error) {
	err := db.Save(setting).Error
	if err != nil {
		return nil, errors.Join(errors.New("Setting update"), err)
	}
	return setting, nil
}

func (setting *Setting) Read() (*[]Setting, error) {
	var settings []Setting
	err := db.Find(&settings).Error
	if err != nil {
		return nil, errors.Join(errors.New("Setting read"), err)
	}
	return &settings, nil
}

func (setting *Setting) Delete() (*Setting, error) {
	err := db.Delete(setting).Error
	if err != nil {
		return nil, errors.Join(errors.New("Setting delete"), err)
	}
	return setting, nil
}
