package security

import "log"

func LoginDash(data Config) error {
	err := dashCheck(data.IP, data.Pwd)
	if err != nil {
		log.Println("Login Failed")
		return err
	} else {
		log.Println("Login Success")
		return nil
	}
}

// Type: gen, sum, img, web, rand
// API: openai, alibaba, deepseek, otherapi, github, gitee
// Target: prompt, repositories
func NewRequest(data Config) error {
	var err error
	err = frequencyCheck(data.IP)
	if err != nil {
		return err
	}
	err = sourceCheck(data.Domain)
	if err != nil {
		return err
	}
	if data.Type == "gen" || data.Type == "sum" || data.Type == "img" {
		err = modelCheck(data.Type, data.API)
		if err != nil {
			return err
		}
	}
	if data.Type == "rand" {
		err = repoCheck(data.API, data.User, data.Repo)
		if err != nil {
			return err
		}
	}
	if data.Type == "gen" {
		err = txtGenCheck(data.Target)
		if err != nil {
			return err
		}
	}
	if data.Type == "img" {
		err = imgGenCheck()
		if err != nil {
			return err
		}
	}
	return nil
}
