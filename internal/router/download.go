package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"urlAPI/internal/data"
	"urlAPI/internal/file"
	"urlAPI/internal/plugin"
	"urlAPI/internal/security"
)

func download() {
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
