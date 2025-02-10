package web

import (
	"backend/cmd/img"
	"backend/internal/file"
	"embed"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"image/png"
	"log"
)

//go:embed fork_icon.png star_icon.png
var icon embed.FS

func arrange(Str string, Width int) []string {
	var maxlen int
	Content := []rune(Str)
	if len(Str) == len(Content) {
		maxlen = (Width - 60) / 15
	} else {
		maxlen = (Width - 60) / 32
	}
	var ret []string
	for i := 0; true; i += maxlen {
		if i+maxlen >= len(Content) {
			ret = append(ret, string(Content[i:len(Content)]))
			break
		} else {
			ret = append(ret, string(Content[i:i+maxlen]))
		}
	}
	return ret
}

func drawRoundedRect(img *image.RGBA) {
	radius := 90
	rect := img.Bounds()

	corners := []image.Point{
		{rect.Min.X + radius, rect.Min.Y + radius},
		{rect.Max.X - radius - 1, rect.Min.Y + radius},
		{rect.Min.X + radius, rect.Max.Y - radius - 1},
		{rect.Max.X - radius - 1, rect.Max.Y - radius - 1},
	}

	draw.Draw(img, image.Rect(rect.Min.X+radius, rect.Min.Y, rect.Max.X-radius, rect.Max.Y), &image.Uniform{image.White}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(rect.Min.X, rect.Min.Y+radius, rect.Max.X, rect.Max.Y-radius), &image.Uniform{image.White}, image.Point{}, draw.Src)

	for _, center := range corners {
		for y := -radius; y <= radius; y++ {
			for x := -radius; x <= radius; x++ {
				if x*x+y*y <= radius*radius {
					img.Set(center.X+x, center.Y+y, image.White)
				}
			}
		}
	}
}

func DrawRepo(logo image.Image, Name, Author, Description, Star, Fork, UUID string) error {
	reader, err := img.SmileySans.ReadFile("ssfonts.ttf")
	if err != nil {
		log.Println("Read font file error")
		return err
	}
	font, err := freetype.ParseFont(reader)
	if err != nil {
		log.Println("Parse font error")
		return err
	}

	starIO, _ := icon.Open("star_icon.png")
	forkIO, _ := icon.Open("fork_icon.png")
	starIcon, err := png.Decode(starIO)
	forkIcon, err := png.Decode(forkIO)
	if err != nil {
		log.Println("Decode icon error")
		return err
	}

	nameLen := len(Name) * 45
	authorLen := len(Author) * 27
	starLen := len(Star) * 27
	forkLen := len(Fork) * 27
	width := max(nameLen, authorLen) + max(starLen, forkLen) + 500

	desriptionContent := arrange(Description, width)
	height := len(desriptionContent)*50 + 300

	templateImg := image.NewRGBA(image.Rect(0, 0, width, height))
	drawRoundedRect(templateImg)
	draw.Draw(templateImg, image.Rect(30, 30, width, height), logo, logo.Bounds().Min, draw.Over)

	drawer := freetype.NewContext()
	drawer.SetDPI(144)
	drawer.SetFont(font)
	drawer.SetDst(templateImg)
	drawer.SetClip(templateImg.Bounds())
	drawer.SetSrc(image.Black)

	drawer.SetFontSize(50)
	drawer.DrawString(Name, freetype.Pt(260, 100))

	drawer.SetFontSize(30)
	drawer.DrawString("by "+Author, freetype.Pt(260, 200))

	draw.Draw(templateImg, image.Rect(width-max(starLen, forkLen)-150, 30, width, height), starIcon, starIcon.Bounds().Min, draw.Over)
	draw.Draw(templateImg, image.Rect(width-max(starLen, forkLen)-150, 140, width, height), forkIcon, forkIcon.Bounds().Min, draw.Over)
	drawer.SetFontSize(30)
	drawer.DrawString(Star, freetype.Pt(width-max(starLen, forkLen)-50, 100))
	drawer.DrawString(Fork, freetype.Pt(width-max(starLen, forkLen)-50, 200))

	drawer.SetFontSize(20)
	for index, content := range desriptionContent {
		drawer.DrawString(content, freetype.Pt(30, 300+index*50))
	}

	return file.Add(file.FileConfig(
		file.WithType("img.save"),
		file.WithUUID(UUID),
		file.WithImg(templateImg)))
}
