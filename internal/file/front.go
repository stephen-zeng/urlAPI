package file

import (
	"errors"
	"log"
)

func Add(data Config) error {
	var err error
	switch data.Type {
	case "img.download":
		err = downloadImg(data.UUID, data.URL)
	case "img.save":
		err = saveImg(data.Img, data.UUID)
	default:
		err = errors.New("unknown file type")
	}
	return err
}

func Del(data Config) error {
	if data.Type == "img" {
		err := delImg(data.UUID)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else if data.Type == "md" {
		return nil
	} else {
		log.Printf("Unknown file type")
		return errors.New("Unknown file type")
	}
}

func Fetch(data Config) ([]byte, error) {
	if data.Type == "img" {
		img, err := fetchImg(data.UUID)
		if err != nil {
			return nil, err
		} else {
			return img, nil
		}
	} else if data.Type == "md" {
		return nil, nil
	} else {
		log.Printf("Unknown file type")
		return nil, errors.New("Unknown file type")
	}
}
