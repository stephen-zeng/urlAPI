package web

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"urlAPI/cmd/img"
	"urlAPI/internal/client"
	"urlAPI/internal/file"
)

var font *truetype.Font
var drawer *freetype.Context

func init() {
	reader, err := img.SmileySans.ReadFile("ssfonts.ttf")
	if err != nil {
		log.Println("Read font file error")
	}
	font, _ = freetype.ParseFont(reader)
	if err != nil {
		log.Println("Parse font error")
	}
	drawer = freetype.NewContext()
	drawer.SetDPI(144)
	drawer.SetFont(font)
	drawer.SetSrc(image.Black)
}

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
			if (Content[i+maxlen] >= 'a' && Content[i+maxlen] <= 'z') || (Content[i+maxlen] >= 'A' && Content[i+maxlen] <= 'Z') {
				ret[len(ret)-1] += "-"
			}
		}
	}
	return ret
}

func drawRoundedRect(img *image.RGBA, option string) {
	radius := 45
	rect := img.Bounds()

	corners := []image.Point{
		{rect.Min.X + radius, rect.Min.Y + radius},
		{rect.Max.X - radius - 1, rect.Min.Y + radius},
		{rect.Min.X + radius, rect.Max.Y - radius - 1},
		{rect.Max.X - radius - 1, rect.Max.Y - radius - 1},
	}

	if option == "fill" {
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
	} else {
		for x := 0; x < radius; x++ {
			for y := 0; y < radius; y++ {
				if (x-radius)*(x-radius)+(y-radius)*(y-radius) > radius*radius {
					img.Set(x, y, image.White)
					img.Set(rect.Dx()-x, y, image.White)
					img.Set(x, rect.Dy()-y, image.White)
					img.Set(rect.Dx()-x, rect.Dy()-y, image.White)
				}
			}
		}
	}
}

func drawRepo(logo image.Image, Name, Author, Description, Star, Fork, UUID string) error {
	starIO, _ := file.IconFS.Open("assets/icon/star_icon.png")
	forkIO, _ := file.IconFS.Open("assets/icon/fork_icon.png")
	starIcon, err := png.Decode(starIO)
	forkIcon, err := png.Decode(forkIO)
	if err != nil {
		log.Println("Decode icon error")
		return err
	}

	var nameLen int
	if len(Name) == len([]rune(Name)) {
		nameLen = len(Name) * 45
	} else {
		nameLen = len(Name) * 80
	}
	Author = "by " + Author
	authorLen := len(Author) * 27
	starLen := len(Star) * 27
	forkLen := len(Fork) * 27
	width := max(nameLen, authorLen) + max(starLen, forkLen) + 500

	desriptionContent := arrange(Description, width)
	height := len(desriptionContent)*50 + 300

	templateImg := image.NewRGBA(image.Rect(0, 0, width, height))
	drawRoundedRect(templateImg, "fill")
	draw.Draw(templateImg, image.Rect(30, 30, width, height), logo, logo.Bounds().Min, draw.Over)

	drawer.SetDst(templateImg)
	drawer.SetClip(templateImg.Bounds())

	drawer.SetFontSize(50)
	drawer.DrawString(Name, freetype.Pt(260, 100))

	drawer.SetFontSize(30)
	drawer.DrawString(Author, freetype.Pt(260, 200))

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

func drawVideo(CoverURL, Name, Author, Description, View, Favorite, Like, Coin, UUID string) error {
	req, err := http.NewRequest("GET", CoverURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	pic, err := jpeg.Decode(resp.Body)
	if err != nil {
		log.Println("Decode jpeg image error")
		return err
	}
	pic = resize.Resize(0, 450, pic, resize.Lanczos3)

	Author = "by " + Author
	var nameLen int
	if len(Name) == len([]rune(Name)) {
		nameLen = len(Name) * 45
	} else {
		nameLen = len([]rune(Name)) * 80
	}
	authorLen := len([]rune(Author)) * 27
	statLen := (max(len(View), len(Like))+max(len(Favorite), len(Coin)))*27 + 250

	templatePic := image.NewRGBA(pic.Bounds())
	draw.Draw(templatePic, templatePic.Bounds(), pic, pic.Bounds().Min, draw.Over)
	drawRoundedRect(templatePic, "boarder")

	width := max(nameLen, authorLen, statLen) + templatePic.Bounds().Dx() + 100
	desriptionContent := arrange(Description, width)
	height := len(desriptionContent)*50 + templatePic.Bounds().Dy() + 100

	templateImg := image.NewRGBA(image.Rect(0, 0, width, height))
	drawRoundedRect(templateImg, "fill")
	draw.Draw(templateImg, image.Rect(30, 30, width, height), templatePic, templatePic.Bounds().Min, draw.Over)

	drawer.SetDst(templateImg)
	drawer.SetClip(templateImg.Bounds())

	drawer.SetFontSize(50)
	drawer.DrawString(Name, freetype.Pt(templatePic.Bounds().Dx()+100, 150))
	drawer.SetFontSize(30)
	drawer.DrawString(Author, freetype.Pt(templatePic.Bounds().Dx()+100, 250))

	drawer.SetFontSize(20)
	for index, content := range desriptionContent {
		drawer.DrawString(content, freetype.Pt(30, templatePic.Bounds().Dy()+index*50+100))
	}

	likeIO, _ := file.IconFS.Open("assets/icon/like_icon.png")
	favIO, _ := file.IconFS.Open("assets/icon/fav_icon.png")
	playIO, _ := file.IconFS.Open("assets/icon/play_icon.png")
	coinIO, _ := file.IconFS.Open("assets/icon/coin_icon.png")
	likeIcon, err := png.Decode(likeIO)
	favIcon, err := png.Decode(favIO)
	playIcon, err := png.Decode(playIO)
	coinIcon, err := png.Decode(coinIO)
	if err != nil {
		log.Println("Decode icon error")
		return err
	}

	draw.Draw(templateImg, image.Rect(templatePic.Bounds().Dx()+100, 300, width, height), playIcon, playIcon.Bounds().Min, draw.Over)
	drawer.DrawString(View, freetype.Pt(templatePic.Bounds().Dx()+180, 350))
	draw.Draw(templateImg, image.Rect(templatePic.Bounds().Dx()+max(len(View), len(Like))*27+200, 300, width, height), favIcon, favIcon.Bounds().Min, draw.Over)
	drawer.DrawString(Favorite, freetype.Pt(templatePic.Bounds().Dx()+max(len(View), len(Like))*27+280, 350))
	draw.Draw(templateImg, image.Rect(templatePic.Bounds().Dx()+100, 400, width, height), likeIcon, likeIcon.Bounds().Min, draw.Over)
	drawer.DrawString(Like, freetype.Pt(templatePic.Bounds().Dx()+180, 450))
	draw.Draw(templateImg, image.Rect(templatePic.Bounds().Dx()+max(len(View), len(Like))*27+200, 400, width, height), coinIcon, coinIcon.Bounds().Min, draw.Over)
	drawer.DrawString(Coin, freetype.Pt(templatePic.Bounds().Dx()+max(len(View), len(Like))*27+280, 450))

	return file.Add(file.FileConfig(
		file.WithType("img.save"),
		file.WithUUID(UUID),
		file.WithImg(templateImg)))
}

func DrawArticle(logo image.Image, ID, Title, Author, Description, UUID, Time string) error {
	titleLen := len(Title) * 25
	var secondTitle string
	if Author != "" {
		secondTitle = "By " + Author
	} else {
		secondTitle = "Time: " + Time
	}
	secondLen := len(secondTitle) * 16
	width := max(titleLen, secondLen) + 60 + logo.Bounds().Dx()
	discriptionContent := arrange(Description, width)
	height := len(discriptionContent)*50 + logo.Bounds().Dy() + 100

	templateImg := image.NewRGBA(image.Rect(0, 0, width, height))
	drawRoundedRect(templateImg, "fill")
	draw.Draw(templateImg, image.Rect(30, 30, width, height), logo, logo.Bounds().Min, draw.Over)

	drawer.SetDst(templateImg)
	drawer.SetClip(templateImg.Bounds())

	drawer.SetFontSize(15)
	drawer.DrawString(ID, freetype.Pt(60+logo.Bounds().Dx(), 50))

	drawer.SetFontSize(32)
	drawer.DrawString(Title, freetype.Pt(60+logo.Bounds().Dx(), 130))

	drawer.SetFontSize(20)
	drawer.DrawString(secondTitle, freetype.Pt(60+logo.Bounds().Dx(), 200))

	drawer.SetFontSize(20)
	for index, content := range discriptionContent {
		drawer.DrawString(content, freetype.Pt(30, 300+index*50))
	}

	return file.Add(file.FileConfig(
		file.WithType("img.save"),
		file.WithUUID(UUID),
		file.WithImg(templateImg)))
}
