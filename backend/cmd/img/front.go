package img

import (
	"encoding/json"
)

func Request(Format, API, Model, Target, Size, IP, Domain, From string) (string, error) {
	retOrigin, err := genRequest(IP, Domain, Model, API, Target, Size, From)
	if err != nil {
		return "", err
	}
	if Format == "json" {
		return retOrigin, nil
	} else {
		response := make(map[string]interface{})
		err = json.Unmarshal([]byte(retOrigin), &response)
		if err != nil {
			return "", err
		}
		return response["url"].(string), nil
	}
}
