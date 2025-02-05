package data

import (
	"encoding/json"
	"log"
)

type Setting struct {
	// 是可导出的类型才能使用GORM
	Name  string `json:"name" gorm:"primaryKey"`
	Value string `json:"value"`
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
