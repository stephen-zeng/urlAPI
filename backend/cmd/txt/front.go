package txt

import (
	"encoding/json"
	"log"
)

func Request(Format, API, Model, Target, Type, IP, Domain string) (string, error) {
	var retOrigin string
	if Type == "sum" {
		retOrigin = "Summary"
		return retOrigin, nil
	} else {
		retOrigin, err := genRequest(IP, Domain, Model, API, Target)
		if err != nil {
			log.Println(err)
			return "", err
		}
		if Format == "json" {
			log.Println(err)
			return retOrigin, nil
		} else {
			ret := make(map[string]interface{})
			err = json.Unmarshal([]byte(retOrigin), &ret)
			if err != nil {
				log.Println(err)
				return "", err
			}
			return ret["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string), nil
		}
	}
}
