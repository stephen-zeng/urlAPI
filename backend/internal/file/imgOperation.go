package file

import (
	"encoding/base64"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
)

var (
	filePath string = "assets/img/"
)

func init() {
	os.MkdirAll(filePath, 0755)
}

func add(uuid, data string) error {
	reader := base64.NewDecoder(base64.StdEncoding,
		strings.NewReader(data))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Println(err)
		return err
	}
	output, err := os.Create(filePath + uuid + ".png")
	if err != nil {
		log.Println(err)
		return err
	}
	err = png.Encode(output, img)
	if err != nil {
		log.Println(err)
		return err
	}
	defer output.Close()
	return nil
}

func del(uuid string) error {
	err := os.Remove(filePath + uuid + ".png")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func fetch(uuid string) (string, error) {
	_, err := os.Stat(filePath + uuid + ".png")
	if err != nil {
		return "", err
	}
	return filePath + uuid + ".png", nil
}
