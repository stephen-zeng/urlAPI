package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
)

var r *gin.Engine

func Start(Port string) {
	gin.SetMode(gin.ReleaseMode)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST"}
	r = gin.New()
	r.Use(cors.New(config))
	setAPI()
	setSession()
	setDash()
	r.Run(":" + Port)
}

var GlobalHTTPClient *http.Client

func init() {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	GlobalHTTPClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}
