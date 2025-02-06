package data

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type Setting struct {
	// 是可导出的类型才能使用GORM
	Name  string `json:"name" gorm:"primaryKey"`
	Value string `json:"value"`
}

func editSetting(name []string, value [][]string, Skip bool) error {
	for index, nameItem := range name {
		valueItem := value[index]
		jsonData, _ := json.Marshal(valueItem)
		newSetting := Setting{
			Name:  nameItem,
			Value: string(jsonData),
		}
		var err *gorm.DB
		if Skip {
			err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&newSetting)
		} else {
			err = db.Save(&newSetting)
		}
		if err.Error != nil {
			log.Println(err.Error)
			return err.Error
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
			log.Println(err.Error)
			return nil, err.Error
		}
		var tmp []string
		_ = json.Unmarshal([]byte(fetch[0].Value), &tmp)
		ret = append(ret, tmp)
	}
	return ret, nil
}
