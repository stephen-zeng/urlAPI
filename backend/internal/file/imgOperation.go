package file

import (
	"encoding/base64"
	"errors"
	"fmt"
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
	fmt.Println("Successfully created folder")
}

func add(uuid, data string) error {
	reader := base64.NewDecoder(base64.StdEncoding,
		strings.NewReader(data))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
		fmt.Println("New File Error")
		return errors.New("file.add.error")
	}
	output, err := os.Create(filePath + uuid + ".png")
	if err != nil {
		log.Fatal(err)
		fmt.Println("New File Error")
		return errors.New("file.add.error")
	}
	err = png.Encode(output, img)
	if err != nil {
		log.Fatal(err)
		fmt.Println("New File Error")
		return errors.New("file.add.error")
	}
	defer output.Close()
	return nil
}

func del(uuid string) error {
	err := os.Remove(filePath + uuid + ".png")
	if err != nil {
		log.Fatal(err)
		fmt.Println("Delete File Error")
		return errors.New("file.del.error")
	}
	return nil
}

func fetch(uuid string) (string, error) {
	_, err := os.Stat(filePath + uuid + ".png")
	if err != nil {
		return "", errors.New("file.fetch.error")
	}
	return filePath + uuid + ".png", nil
}
