package security

import (
	"errors"
	"log"
)

// 切片不能是常量
var txtAPI = []string{"openai", "alibaba", "deepseek", "otherapi"}
var imgAPI = []string{"openai", "alibaba"}

func modelCheck(Type, API string) error {
	if Type == "txt.gen" {
		for _, api := range txtAPI {
			if API == api {
				return nil
			}
		}
		log.Println("The request " + API + " is NOT valid.")
		return errors.New("invalid request")
	} else if Type == "img.gen" {
		for _, api := range imgAPI {
			if API == api {
				return nil
			}
		}
		log.Println("The request " + API + " is NOT valid.")
		return errors.New("invalid request")
	} else {
		log.Println("The request " + API + " is NOT valid.")
		return errors.New("invalid request")
	}
}
