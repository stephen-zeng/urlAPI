package file

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
)

var (
	filePath string = "assets/tmp"
)

func add(data map[string]interface{}) interface{} {
	reader := base64.NewDecoder(base64.StdEncoding,
		strings.NewReader(data["data"].(string)))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
		fmt.Println("New File Error")
		return "file.add.error"
	}
	output, err := os.Create(filePath + data["uuid"].(string) + ".png")
	if err != nil {
		log.Fatal(err)
		fmt.Println("New File Error")
		return "file.add.error"
	}
	defer output.Close()
	err = png.Encode(output, img)
	if err != nil {
		log.Fatal(err)
		fmt.Println("New File Error")
		return "file.add.error"
	}
	return "file.add.success"
}

func del(data map[string]interface{}) interface{} {
	err := os.Remove(filePath + data["uuid"].(string) + ".png")
	if err != nil {
		log.Fatal(err)
		fmt.Println("Delete File Error")
		return "file.del.error"
	}
	return "file.del.success"
}

func fetch(data map[string]interface{}) interface{} {
	_, err := os.Stat(filePath + data["uuid"].(string) + ".png")
	if err != nil {
		return "file.fetch.error"
	}
	return filePath + data["uuid"].(string) + ".png"
}
