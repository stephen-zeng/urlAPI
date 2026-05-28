package op

import (
	"github.com/pkg/errors"
	"urlAPI/internal/model"
)

func HandleSession(request Session, authSession model.Session) (Session, error) {
	response := request
	if err := login(&response, &authSession); err != nil {
		return response, errors.WithStack(err)
	}

	var err error
	switch response.Operation {
	case "login":
	case "logout":
		err = logout(&authSession)
	case "exit":
		err = exit(&authSession)
	case "newRepo":
		err = newRepo(&response)
	case "refreshRepo":
		err = refreshRepo(&response)
	case "delRepo":
		err = delRepo(&response)
	case "fetchRepo":
		err = fetchRepo(&response)
	case "fetchTask":
		err = fetchTask(&response)
	case "fetchTaskStats":
		err = fetchTaskStats(&response)
	case "fetchSettings":
		err = fetchSettings(&response)
	case "editSettings":
		err = editSettings(&response)
	}
	if err != nil {
		return response, errors.WithStack(err)
	}
	return response, nil
}
