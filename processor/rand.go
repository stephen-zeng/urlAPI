package processor

import (
	"github.com/pkg/errors"
	"math/rand"
	"urlAPI/database"
)

func (info *Rand) Process(data *database.Task) error {
	info.Return = database.SettingMap["rand"][2]
	if info.API == "" {
		info.API = database.SettingMap["rand"][3]
	}
	var content []string
	var ok bool
	content, ok = database.RepoMap[info.API+";"+info.Target]
	if !ok {
		data.Status = "failed"
		data.Return = "Repo not found"
		return errors.WithStack(errors.New("Process Rand Repo not found"))
	}
	length := len(content)
	index := rand.Intn(length)
	info.Return = content[index]
	data.Return = info.Return
	data.Status = "success"
	return nil
}
