package file

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	imgPath string = "assets/img/"
)

func init() {
	os.MkdirAll(imgPath, 0755)
}

func addImg(UUID, URL string) error {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	reader := bytes.NewReader(content)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Println(err)
		return err
	}
	output, err := os.Create(imgPath + UUID + ".png")
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

func delImg(uuid string) error {
	err := os.Remove(imgPath + uuid + ".png")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func fetchImg(UUID string) ([]byte, error) {
	img, err := os.ReadFile(imgPath + UUID + ".png")
	if err != nil {
		log.Println("ReadFile error")
		return nil, err
	}
	return img, nil
}
