package file

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 这里给的裸的response
func Add(uuid, source, data string) error {
	var url string
	var response map[string]interface{}
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		fmt.Println("JSON Unmarshal Error")
		return errors.New("file.add.error")
	}
	if source == "aliyun" {
		if response["output"].(map[string]interface{})["task_status"].(string) != "SUCCEEDED" {
			return errors.New("file.add.error")
		} else {
			url = response["output"].(map[string]interface{})["results"].([]interface{})[0].(map[string]interface{})["url"].(string)
		}
	} else if source == "openai" {
		url = response["data"].([]interface{})[0].(map[string]interface{})["url"].(string)
	}
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get(url)
	if err != nil {
		return errors.New("file.add.error")
	}
	defer resp.Body.Close()
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("file.add.error")
	}
	b64Img := base64.StdEncoding.EncodeToString(image)
	err = add(uuid, b64Img)
	if err != nil {
		fmt.Println("file.add.error")
		return errors.New("file.add.error")
	} else {
		return nil
	}
}

func Del(uuid string) error {
	err := del(uuid)
	if err != nil {
		return errors.New("file.del.error")
	} else {
		return nil
	}
}

func Fetch(uuid string) (string, error) {
	ret, err := fetch(uuid)
	if err != nil {
		return "", errors.New("file.fetch.error")
	} else {
		return ret, nil
	}
}
