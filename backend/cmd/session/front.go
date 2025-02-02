package session

import (
	"backend/internal/data"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

// string --> pwd, token
func Auth(dat Config) (string, error) {
	dash, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"dash", "dashallowedip"})))
	if err != nil {
		return "", err
	}
	flag := false
	for _, item := range dash[1] {
		rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(item), `\*`, ".*") + "$"
		match, err := regexp.MatchString(rgx, dat.IP)
		if err != nil {
			continue
		}
		if match {
			flag = true
			break
		}
	}
	if !flag {
		log.Println("Authentication Failed by IP")
		return "", errors.New("Authentication Failed by IP")
	}
	if dat.Token == dash[0][0] {
		return "pwd", nil
	}
	tokens, err := data.FetchSession(data.DataConfig(data.WithToken(dat.Token)))
	if err != nil {
		return "", errors.New("Authentication Failed")
	}
	exp := tokens[0].Expire
	if time.Now().After(exp) {
		data.DelSession(data.DataConfig(data.WithToken(dat.Token)))
		return "", errors.New("Token expired")
	} else {
		return "token", nil
	}
}

func New(dat Config) (SessionResponse, error) {
	var ret SessionResponse
	var err error
	if dat.Operation == "login" {
		if dat.Type == "token" {
			ret.Token = dat.Token
			return ret, nil
		} else {
			ret.Token, err = login(SessionConfig(WithTerm(dat.Term)))
			if err != nil {
				return SessionResponse{}, err
			} else {
				return ret, nil
			}
		}
	}
	if dat.Operation == "logout" {
		return SessionResponse{}, logout(SessionConfig(WithToken(dat.Token)))
	}
	if dat.Operation == "exit" {
		return SessionResponse{}, exit(SessionConfig(WithToken(dat.Token)))
	}
	if dat.Operation == "clear" {
		return SessionResponse{}, data.InitSession(data.DataConfig(data.WithType("restore")))
	}
	if dat.Operation == "fetch" {
		response, err := fetch(SessionConfig(WithPart(dat.Part)))
		if err != nil {
			return SessionResponse{}, err
		} else {
			ret.Token = dat.Token
			ret.Part = dat.Part
			ret.Name = response.Name
			ret.Setting = response.Setting
			return ret, nil
		}
	}
	if dat.Operation == "edit" {
		response, err := edit(SessionConfig(
			WithPart(dat.Part),
			WithEdit(dat.Edit)))
		if err != nil {
			return SessionResponse{}, err
		} else {
			ret.Token = dat.Token
			ret.Part = dat.Part
			ret.Name = response.Name
			return ret, nil
		}
	}
	if dat.Operation == "task" {
		response, err := data.FetchTask(data.DataConfig(data.WithUUID(dat.By)))
		if err != nil {
			return SessionResponse{}, err
		} else {
			ret.Task = response
			return ret, nil
		}
	}
	return SessionResponse{}, errors.New("Invalid Session Operation")
}
