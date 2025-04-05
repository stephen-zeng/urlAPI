package processor

import (
	"github.com/pkg/errors"
	"urlAPI/database"
)

func (info *Session) Process(data *database.Session) error {
	var err error
	if err = login(info, data); err != nil {
		return errors.WithStack(err)
	}
	switch info.Operation {
	case "logout":
		err = logout(info, data)
	case "exit":
		err = exit(info, data)
	case "newRepo":
		err = newRepo(info, data)
	case "refreshRepo":
		err = refreshRepo(info, data)
	case "delRepo":
		err = delRepo(info, data)
	case "fetchRepo":
		err = fetchRepo(info, data)
	case "fetchTask":
		err = fetchTask(info, data)
	case "fetchSetting":
		err = fetchSetting(info, data)
	case "editSetting":
		err = editSetting(info, data)
	case "fetchPrompt":
		err = fetchPrompt(info, data)
	case "editPrompt":
		err = editPrompt(info, data)
	case "newPrompt":
		err = newPrompt(info, data)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
