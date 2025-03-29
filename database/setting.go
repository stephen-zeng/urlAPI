package database

import (
	"encoding/json"
	"errors"
	"io"
	"urlAPI/file"
)

func settingInit() error {
	f, _ := file.Settings.Open("setting.json")
	d, _ := io.ReadAll(f)
	var settingsInit SettingInit
	_ = json.Unmarshal(d, &settingsInit)
	for index, settingInitName := range settingsInit.Names {
		settingInitList := settingsInit.Edits[index]
		dbSettingWriter := Setting{
			Name: settingInitName,
		}
		dbSettingList := SettingMap[settingInitName]
		if len(dbSettingList) < len(settingInitList) {
			dbSettingList = append(dbSettingList, settingInitList[len(dbSettingList):]...)
		}
		jsonList, _ := json.Marshal(dbSettingList)
		dbSettingWriter.Value = string(jsonList)
		_ = dbSettingWriter.Update()
	}
	return nil
}

func (setting *Setting) Create() error {
	err := db.Create(&setting).Error
	if err != nil {
		return errors.Join(errors.New("Setting create"), err)
	}
	var tmp []string
	err = json.Unmarshal([]byte(setting.Value), &tmp)
	if err != nil {
		return errors.Join(errors.New("Setting create"), err)
	}
	SettingMap[setting.Name] = tmp
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
	err := db.Where("name=?", setting.Name).Find(&settings).Error
	if len(settings) == 0 {
		err = errors.New("Setting not found")
	}
	ret := DBList{
		SettingList: settings,
	}
	if err != nil {
		return &ret, errors.Join(errors.New("Setting read"), err)
	}
	return &ret, nil
}

func (setting *Setting) Delete() error {
	err := db.Delete(setting).Error
	if err != nil {
		return errors.Join(errors.New("Setting delete"), err)
	}
	return nil
}
