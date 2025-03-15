package session

import (
	"errors"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"
	"urlAPI/internal/data"
)

func Auth(dat Config) (string, error) {
	dash, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"dash", "dashallowedip"})))
	if err != nil {
		return "", err
	}
	flag := false
	for _, item := range dash[1] {
		if dat.SessionIP == "::1" || dat.SessionIP == "127.0.0.1" {
			flag = true
			break
		}
		rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(item), `\*`, ".*") + "$"
		match, err := regexp.MatchString(rgx, dat.SessionIP)
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
	if dat.SessionToken == dash[0][0] {
		return "pwd", nil
	}
	tokens, err := data.FetchSession(data.DataConfig(data.WithSessionToken(dat.SessionToken)))
	if err != nil {
		return "", errors.New("Authentication Failed")
	}
	exp := tokens[0].Expire
	if time.Now().After(exp) {
		data.DelSession(data.DataConfig(data.WithSessionToken(dat.SessionToken)))
		return "", errors.New("Token expired")
	} else {
		return "token", nil
	}
}

func New(dat Config) (SessionResponse, error) {
	var ret SessionResponse
	var err error
	ret.SessionToken = dat.SessionToken
	switch dat.Operation {
	case "login":
		if dat.SessionType == "pwd" {
			ret.SessionToken, err = newLogin(dat.SessionTerm)
		}
	case "logout":
		return SessionResponse{}, logout(dat.SessionToken)
	case "exit":
		return SessionResponse{}, exit(dat.SessionToken)
	case "fetchSetting":
		ret.SettingPart = dat.SettingPart
		ret.SettingName, ret.SettingData, err = fetchSetting(dat.SettingPart)
	case "editSetting":
		ret.SettingPart = dat.SettingPart
		ret.SettingName, err = editSetting(dat.SettingPart, dat.SettingEdit)
	case "fetchTask":
		ret.TaskData, err = data.FetchTask(data.DataConfig(
			data.WithType(dat.TaskCatagory),
			data.WithBy(dat.TaskBy),
			data.WithTaskPage(dat.TaskPage)))
		ret.TaskMaxPage = (len(ret.TaskData) / 500) + 1
		sort.Slice(
			ret.TaskData,
			func(i, j int) bool {
				return ret.TaskData[i].Time.After(ret.TaskData[j].Time)
			},
		)
		switch {
		case dat.TaskPage == -1:
			break
		case dat.TaskPage*500 > len(ret.TaskData):
			ret.TaskData = ret.TaskData[(dat.TaskPage-1)*500:]
		case dat.TaskPage*500 <= len(ret.TaskData):
			ret.TaskData = ret.TaskData[(dat.TaskPage-1)*500:]
		}
	case "fetchRepo":
		ret.RepoData, err = data.FetchRepo(data.DataConfig())
	case "newRepo":
		err = data.NewRepo(data.DataConfig(
			data.WithAPI(dat.RepoAPI),
			data.WithRepoInfo(dat.RepoInfo),
		))
	case "refreshRepo":
		err = data.RefreshRepo(data.DataConfig(
			data.WithUUID(dat.RepoUUID),
		))
	case "delRepo":
		err = data.DelRepo(data.DataConfig(
			data.WithUUID(dat.RepoUUID),
		))
	default:
		return SessionResponse{}, errors.New("Invalid Session Operation")
	}
	if err != nil {
		return SessionResponse{}, err
	} else {
		return ret, nil
	}
}
