package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"urlAPI/file"
	"urlAPI/util"
)

func settingInit() error {
	f, _ := file.Settings.Open("setting.json")
	d, _ := io.ReadAll(f)
	var settingsInit SettingInit
	_ = json.Unmarshal(d, &settingsInit)
	for index, settingInitName := range settingsInit.Names {
		settingInitList := settingsInit.Edits[index]
		if len(settingInitList) == 0 {
			settingInitList = append(settingInitList, util.GetShortRandomString(16))
		}
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
	if err := db.Create(&setting).Error; err != nil {
		return errors.WithStack(err)
	}
	var tmp []string
	if err := json.Unmarshal([]byte(setting.Value), &tmp); err != nil {
		return errors.WithStack(err)
	}
	SettingMap[setting.Name] = tmp
	return nil
}

func (setting *Setting) Update() error {
	if err := db.Save(setting).Error; err != nil {
		return errors.WithStack(err)
	}
	var tmp []string
	if err := json.Unmarshal([]byte(setting.Value), &tmp); err != nil {
		return errors.WithStack(err)
	}
	SettingMap[setting.Name] = tmp
	return nil
}

func (setting *Setting) Read() (*DBList, error) {
	var settings []Setting
	err := db.Where("name=?", setting.Name).Find(&settings).Error
	if len(settings) == 0 {
		err = errors.WithStack(errors.New("Setting not found"))
	}
	ret := DBList{
		SettingList: settings,
	}
	return &ret, errors.WithStack(err)
}

func (setting *Setting) Delete() error {
	return errors.WithStack(db.Delete(setting).Error)
}
