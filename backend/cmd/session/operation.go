package session

import "backend/cmd/set"

func fetch(dat Config) (SessionResponse, error) {
	response, err := set.Fetch(set.SetConfig(set.WithPart(dat.Part)))
	if err != nil {
		return SessionResponse{}, err
	} else {
		return SessionResponse{
			Name:    response.Name,
			Setting: response.Setting,
		}, nil
	}
}

func edit(dat Config) (SessionResponse, error) {
	response, err := set.Edit(set.SetConfig(
		set.WithPart(dat.Part),
		set.WithEdit(dat.Edit)))
	if err != nil {
		return SessionResponse{}, err
	} else {
		return SessionResponse{
			Name: response.Name,
		}, nil
	}
}
