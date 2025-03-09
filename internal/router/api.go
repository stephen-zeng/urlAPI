package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"regexp"
	"urlAPI/cmd/img"
	"urlAPI/cmd/txt"
	"urlAPI/cmd/web"
	"urlAPI/internal/data"
	"urlAPI/internal/file"
	"urlAPI/internal/plugin"
	"urlAPI/internal/security"
)

func getScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https://"
	}
	if scheme := c.GetHeader("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http://"
}

func getDeviceType(ua string) string {
	mobileRegexp := `(?i)(Mobile|Tablet|Android|iOS|iPhone|iPad|iPod)`
	desktopRegexp := `(?i)(Desktop|Windows|Macintosh|Linux)`
	botRegexp := `(?i)(Bot)`
	matched, _ := regexp.MatchString(mobileRegexp, ua)
	if matched {
		return "Mobile"
	}
	matched, _ = regexp.MatchString(desktopRegexp, ua)
	if matched {
		return "Desktop"
	}
	matched, _ = regexp.MatchString(botRegexp, ua)
	if matched {
		return "Bot"
	}
	return ""
}

func downloadRequest() {
	r.GET("/download", func(c *gin.Context) {
		var err error
		referer, err := url.Parse(c.Request.Referer())
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		domain := referer.Hostname()
		Img := c.Query("img")
		Md := c.Query("md")
		device := getDeviceType(c.GetHeader("User-Agent"))
		var id string
		if Img != "" {
			id = Img
		} else {
			id = Md
		}
		err = security.NewRequest(security.SecurityConfig(
			security.WithIP(c.ClientIP()),
			security.WithDomain(domain),
			security.WithType("download"),
			security.WithTarget(id)))
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		var dat []byte
		var suffix string
		if Img != "" {
			suffix = ".png"
			dat, err = file.Fetch(file.FileConfig(
				file.WithType("img"),
				file.WithUUID(Img)))
			c.Header("Content-Type", "image/png")
		} else if Md != "" {
			suffix = ".md"
			dat, err = file.Fetch(file.FileConfig(
				file.WithType("md"),
				file.WithUUID(Md)))
			c.Header("Content-Type", "text/plain")
		} else {
			c.JSON(400, gin.H{
				"error": "Unknow file type",
			})
		}
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		region, err := plugin.GetRegion(plugin.PluginConfig(plugin.WithIP(c.ClientIP())))
		if err != nil {
			log.Println("Region fetch failed")
		}
		_, err = data.NewTask(data.DataConfig(
			data.WithTaskIP(c.ClientIP()),
			data.WithTaskRegion(region.Region),
			data.WithType("文件下载"),
			data.WithTaskStatus("success"),
			data.WithTaskTarget(id),
			data.WithTaskReferer(referer.String()),
			data.WithTaskDevice(device),
		))
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		} else {
			c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="download%s"`, suffix))
			c.Header("Accept-Length", fmt.Sprintf("%d", len(dat)))
			c.Writer.Write(dat)
		}
	})
}

func txtRequest() {
	r.GET("/txt", func(c *gin.Context) {
		if c.Request.Referer() == "" {
			c.JSON(301, gin.H{
				"error": "Empty Source",
			})
			return
		}
		fallbackURL := data.FallbackURL
		config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"txt"})))
		if len(config[0]) > 4 {
			fallbackURL = config[0][4]
		}
		referer, err := url.Parse(c.Request.Referer())
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
			return
		}
		format := c.Query("format")
		api := c.Query("api")
		model := c.Query("model")
		prompt := c.Query("prompt")
		device := getDeviceType(c.GetHeader("User-Agent"))
		response, err := txt.GenRequest(
			c.ClientIP(), getScheme(c)+c.Request.Host,
			model, api, prompt, device, referer)
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
		} else {
			if format == "json" {
				c.JSON(200, response)
			} else if format == "txt" {
				c.String(200, response.Response)
			} else {
				c.Redirect(302, response.URL)
			}
		}
	})
}

func imgRequest() {
	r.GET("/img", func(c *gin.Context) {
		if c.Request.Referer() == "" {
			c.JSON(301, gin.H{
				"error": "Empty Source",
			})
			return
		}
		fallbackURL := data.FallbackURL
		config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"img"})))
		if len(config[0]) > 3 {
			fallbackURL = config[0][3]
		}
		referer, err := url.Parse(c.Request.Referer())
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
			return
		}
		format := c.Query("format")
		api := c.Query("api")
		model := c.Query("model")
		prompt := c.Query("prompt")
		size := c.Query("size")
		device := getDeviceType(c.GetHeader("User-Agent"))
		response, err := img.GenRequest(
			c.ClientIP(),
			model, api, prompt, size,
			getScheme(c)+c.Request.Host,
			device, referer)
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
		} else {
			if format == "json" {
				c.JSON(200, response)
			} else {
				c.Redirect(302, response.URL)
			}
		}
	})
}

func webRequest() {
	r.GET("/web", func(c *gin.Context) {
		if c.Request.Referer() == "" {
			c.JSON(301, gin.H{
				"error": "Empty Source",
			})
			return
		}
		fallbackURL := data.FallbackURL
		config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"web"})))
		if len(config[0]) > 4 {
			fallbackURL = config[0][4]
		}
		referer, err := url.Parse(c.Request.Referer())
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
			return
		}
		format := c.Query("format")
		img := c.Query("img")
		device := getDeviceType(c.GetHeader("User-Agent"))
		//sum := c.Query("sum")
		var targetURL *url.URL
		var target string
		var response web.WebResponse
		if img[0] == 'h' {
			target = img
			targetURL, err = url.Parse(img)
			response, err = web.ImgRequest(
				c.ClientIP(),                // IP
				getScheme(c)+c.Request.Host, // https://api.example.com
				targetURL.Hostname(),        // github.com ...
				target, device, referer)
		} else {
			log.Println("Improper request")
			c.Redirect(302, fallbackURL)
			return
		}
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
			return
		}
		if format == "json" {
			c.JSON(200, response)
		} else {
			c.Redirect(302, response.URL)
		}
	})
}

func randRequest() {
	r.GET("/rand", func(c *gin.Context) {
		if c.Request.Referer() == "" {
			c.JSON(301, gin.H{
				"error": "Empty Source",
			})
			return
		}
		fallbackURL := data.FallbackURL
		list, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"rand"})))
		fallbackURL = list[0][2]
		referer, err := url.Parse(c.Request.Referer())
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
			return
		}
		domain := referer.Hostname()
		format := c.Query("format")
		api := c.Query("api")
		user := c.Query("user")
		repo := c.Query("repo")
		device := getDeviceType(c.GetHeader("User-Agent"))
		err = security.NewRequest(security.SecurityConfig(
			security.WithType("rand"),
			security.WithAPI(api),
			security.WithTarget(user+"/"+repo),
			security.WithDomain(domain),
		))
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
			return
		}
		region, err := plugin.GetRegion(plugin.PluginConfig(plugin.WithIP(c.ClientIP())))
		if err != nil {
			log.Println("Region fetch failed")
		}
		id, err := data.NewTask(data.DataConfig(
			data.WithAPI(api),
			data.WithType("随机图片"),
			data.WithTaskTarget(user+"/"+repo),
			data.WithTaskRegion(region.Region),
			data.WithTaskIP(c.ClientIP()),
			data.WithTaskReferer(referer.String()),
			data.WithTaskDevice(device),
		))
		response, err := plugin.Request(plugin.PluginConfig(
			plugin.WithAPI(api),
			plugin.WithRepo(user+"/"+repo),
		))
		if err != nil {
			editErr := data.EditTask(data.DataConfig(
				data.WithUUID(id),
				data.WithTaskStatus("failed"),
				data.WithTaskReturn(err.Error()),
			))
			if editErr != nil {
				err = editErr
			}
			log.Println(err)
			c.Redirect(302, fallbackURL)
			return
		}
		err = data.EditTask(data.DataConfig(
			data.WithUUID(id),
			data.WithTaskStatus("success"),
			data.WithTaskReturn(response.URL),
		))
		if err != nil {
			log.Println(err)
			c.Redirect(302, fallbackURL)
			return
		}
		if format == "json" {
			c.JSON(200, gin.H{
				"URL":  response.URL,
				"API":  api,
				"Repo": user + "/" + repo,
			})
		} else {
			c.Redirect(302, response.URL)
		}
	})
}

func setAPI() {
	txtRequest()
	imgRequest()
	randRequest()
	webRequest()
	downloadRequest()
}
