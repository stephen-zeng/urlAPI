package database

import (
	"encoding/json"
	"errors"
	"io"
	"urlAPI/file"
)

func settingInit() error {
	err := db.AutoMigrate(&Setting{})
	if err != nil {
		return errors.Join(errors.New("Setting Init"), err)
	}
	f, _ := file.Settings.Open("settings.json")
	d, _ := io.ReadAll(f)
	var settingsInit SettingInit
	_ = json.Unmarshal(d, &settingsInit)
	for index, settingInitName := range settingsInit.Names {
		settingInitList := settingsInit.Edits[index]
		dbSettingFetcher := Setting{
			Name: settingInitName,
		}
		dbList, _ := dbSettingFetcher.Read()
		dbSetting := dbList.SettingList[0]
		dbSettingJsonList := dbSetting.Value
		var dbSettingList []string
		_ = json.Unmarshal([]byte(dbSettingJsonList), &dbSettingList)
		if len(dbSettingList) < len(settingInitList) {
			dbSettingList = append(dbSettingList, settingInitList[len(dbSettingJsonList):]...)
		}
		editedSettingJsonList, _ := json.Marshal(dbSettingList)
		dbSetting.Value = string(editedSettingJsonList)
		_ = dbSetting.Update()
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
