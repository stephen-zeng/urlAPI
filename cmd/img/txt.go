package img

import (
	"embed"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"log"
	"unicode/utf8"
	"urlAPI/internal/file"
)

//go:embed ssfonts.ttf
var SmileySans embed.FS

func arrange(Str string) []string {
	Content := []rune(Str)
	var ret []string
	for i := 0; true; i += 20 {
		if i+20 >= len(Content) {
			ret = append(ret, string(Content[i:len(Content)]))
			break
		} else {
			ret = append(ret, string(Content[i:i+20]))
		}
	}
	return ret
}

func TxtDrawRequest(UUID, Str, From string) (ImgResponse, error) {
	reader, err := SmileySans.ReadFile("ssfonts.ttf")
	if err != nil {
		log.Println("Failed to load font file")
		return ImgResponse{}, err
	}
	font, err := freetype.ParseFont(reader)
	if err != nil {
		log.Println("Failed to parse font file")
		return ImgResponse{}, err
	}

	Content := arrange(Str)

	templateImg := image.NewRGBA(image.Rect(0, 0, (25 + 40*utf8.RuneCountInString(Content[0])), (60*len(Content) + 13)))
	drawer := freetype.NewContext()
	drawer.SetDPI(144)
	drawer.SetDst(templateImg)
	drawer.SetClip(templateImg.Bounds())

	drawer.SetFont(font)
	drawer.SetFontSize(25)

	for index, content := range Content {
		drawer.SetSrc(image.NewUniform(color.RGBA{100, 100, 100, 255}))
		drawer.DrawString(content, freetype.Pt(15, 60*(index+1)+2))
		drawer.SetSrc(image.White)
		drawer.DrawString(content, freetype.Pt(13, 60*(index+1)))
	}

	return ImgResponse{
			InitPrompt:   Str,
			ActualPrompt: Str,
			URL:          From + "/download?img=" + UUID,
		}, file.Add(file.FileConfig(
			file.WithType("img.save"),
			file.WithUUID(UUID),
			file.WithImg(templateImg)))
}
