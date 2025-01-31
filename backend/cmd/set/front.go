package set

import (
	"backend/internal/data"
	"encoding/json"
)

func Fetch(data string) (string, error) {
	dat := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &dat)
	if err != nil {
		return "", err
	}
	retStr, err := fetch(SetConfig(WithPart(dat["part"].(string))))
	if err != nil {
		return "", err
	}
	ret, err := json.Marshal(retStr)
	if err != nil {
		return "", err
	} else {
		return string(ret), nil
	}
}

// 蠢死了interface切片不能自己转string切片
func Edit(data string) error {
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(data), &dat)
	if err != nil {
		return err
	}
	var editList [][]string
	for _, itemlist := range dat["edit"].([]interface{}) {
		var tmpList []string
		for _, item := range itemlist.([]interface{}) {
			tmpList = append(tmpList, item.(string))
		}
		editList = append(editList, tmpList)
	}
	err = edit(SetConfig(
		WithPart(dat["part"].(string)),
		WithEdit(editList)))
	return err
}

func RePwd() (string, error) {
	return repwd()
}

func Restore() (string, error) {
	return data.InitSetting(data.DataConfig(data.WithType("restore")))
}
