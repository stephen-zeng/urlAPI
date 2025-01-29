package file

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"time"
)

func AddImg(uuid, source, url string) error {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	b64Img := base64.StdEncoding.EncodeToString(image)
	err = add(uuid, b64Img)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func DelImg(uuid string) error {
	err := del(uuid)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func FetchImg(uuid string) (string, error) {
	ret, err := fetch(uuid)
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		return ret, nil
	}
}
