package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"urlAPI/util"
)

func promptInit() error {
	txtPromptFullList := SettingMap["txtacceptprompt"]
	imgPromptFullList := SettingMap["imgacceptprompt"]
	txtSortedList := util.GetKeyList(&txtPromptFullList)
	imgSortedList := util.GetKeyList(&imgPromptFullList)
	// 写入txt
	for index, txt := range txtSortedList {
		jsonData, _ := json.Marshal(txt)
		dbWriter := Prompt{
			For:     "txt",
			Letter:  index,
			Prompts: string(jsonData),
		}
		if err := dbWriter.Update(); err != nil {
			return errors.WithStack(err)
		}
	}

	//写入img
	for index, img := range imgSortedList {
		jsonData, _ := json.Marshal(img)
		dbWriter := Prompt{
			For:     "img",
			Letter:  index,
			Prompts: string(jsonData),
		}
		if err := dbWriter.Update(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (prompt *Prompt) Update() error {
	return errors.WithStack(db.Where("for=?", prompt.For).Where("letter=?", prompt.Letter).Save(prompt).Error)
}

func (prompt *Prompt) Read() (*DBList, error) {
	var prompts []Prompt
	err := db.Where("letter=?", prompt.Letter).Where("for=?", prompt.For).Find(&prompts).Error
	if len(prompts) == 0 {
		err = gorm.ErrRecordNotFound
	}
	ret := DBList{
		PromptList: prompts,
	}
	return &ret, errors.WithStack(err)
}

func (prompt *Prompt) Delete() error {
	return errors.WithStack(db.Where("for=?", prompt.For).Where("letter=?", prompt.Letter).Delete(prompt).Error)
}

func (prompt *Prompt) Create() error {
	return errors.WithStack(db.Create(&prompt).Error)
}
