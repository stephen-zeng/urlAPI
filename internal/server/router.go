package server

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
	"urlAPI/internal/server/handles"
	"urlAPI/internal/server/middleware"
	"urlAPI/static"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/**
 * @brief 创建并配置 HTTP 路由引擎。
 * @return *gin.Engine 已注册路由、中间件和静态资源的 Gin 引擎。
 */
func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Authorization", "X-API-Key", "X-Provider"}
	config.AllowMethods = []string{"GET", "POST"}

	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(requestLogger())
	r.Use(cors.New(config))
	rootFS, _ := fs.Sub(static.StaticFS, "dist")
	assetsFS, _ := fs.Sub(static.StaticFS, "dist/assets")
	r.StaticFS("/assets", http.FS(assetsFS))
	tpl := template.Must(template.ParseFS(rootFS, "*.html"))
	r.SetHTMLTemplate(tpl)

	r.GET("/txt", middleware.GeneralSecurityMiddleware("txt"), middleware.TextSecurityMiddleware(), handles.TxtHandler)
	r.GET("/img", middleware.GeneralSecurityMiddleware("img"), middleware.ImageSecurityMiddleware(), handles.ImgHandler)
	r.GET("/rand", middleware.GeneralSecurityMiddleware("rand"), middleware.RandomSecurityMiddleware(), handles.RandHandler)
	r.GET("/web", middleware.GeneralSecurityMiddleware("web"), middleware.WebSecurityMiddleware(), handles.WebHandler)
	r.GET("/download", middleware.GeneralSecurityMiddleware("download"), middleware.DownloadSecurityMiddleware(), handles.DownloadHandler)
	r.POST("/session", handles.SessionHandler)

	// OpenAI-compatible API v1 (requires API Key)
	v1 := r.Group("/v1")
	v1.Use(middleware.APIKeyAuthMiddleware(middleware.AuthConfig{Mode: middleware.AuthModeRequired}))
	{
		v1.POST("/responses", handles.ResponsesHandler)
		v1.POST("/chat/completions", handles.ChatCompletionHandler)
		v1.POST("/embeddings", handles.EmbeddingsHandler)
		v1.GET("/models", handles.ModelsHandler)
	}

	// API Key management (requires admin session)
	admin := r.Group("/admin/apikeys")
	{
		admin.GET("", handles.ListAPIKeysHandler)
		admin.POST("", handles.CreateAPIKeyHandler)
		admin.DELETE("/:id", handles.DeleteAPIKeyHandler)
		admin.PATCH("/:id", handles.UpdateAPIKeyHandler)
	}

	r.NoMethod(handles.MethodNotAllowedHandler)
	r.NoRoute(handles.StaticHandler)
	return r
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("[HTTP] --> %s %s | IP: %s | UA: %s", c.Request.Method, c.Request.URL.RequestURI(), c.ClientIP(), c.Request.UserAgent())
		c.Next()
		log.Printf("[HTTP] <-- %s %s | Status: %d | Duration: %s", c.Request.Method, c.Request.URL.RequestURI(), c.Writer.Status(), time.Since(start))
	}
}

/**
 * @brief 启动 HTTP 服务。
 * @param port 监听端口。
 * @return error 服务启动失败时返回错误。
 */
func Run(port string) error {
	log.Printf("The server will be running on port %s", port)
	return NewRouter().Run(":" + port)
}
