package file

import (
	"errors"
	"log"
)

func Add(data Config) error {
	if data.Type == "img" {
		err := addImg(data.UUID, data.URL)
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
