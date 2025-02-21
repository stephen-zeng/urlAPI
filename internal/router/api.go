package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"urlAPI/cmd/img"
	"urlAPI/cmd/txt"
	"urlAPI/cmd/web"
	"urlAPI/internal/data"
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
		response, err := txt.GenRequest(
			c.ClientIP(), getScheme(c)+c.Request.Host,
			model, api, prompt, referer)
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
		response, err := img.GenRequest(
			c.ClientIP(),
			model, api, prompt, size,
			getScheme(c)+c.Request.Host,
			referer)
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
		//sum := c.Query("sum")
		var targetURL *url.URL
		var target string
		var response web.WebResponse
		if img != "" {
			target = img
			targetURL, err = url.Parse(img)
			response, err = web.ImgRequest(
				c.ClientIP(),                // IP
				getScheme(c)+c.Request.Host, // https://api.example.com
				targetURL.Hostname(),        // github.com ...
				target, referer)
		} else {
			log.Println("Empty request")
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
}
