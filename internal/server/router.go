package server

import (
	"context"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlAPI/internal/server/handles"
	"urlAPI/internal/server/middleware"
	"urlAPI/static"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST"}

	r := gin.New()
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
	r.NoRoute(handles.StaticHandler)
	return r
}

// shutdownTimeout bounds how long in-flight requests may run after a
// termination signal.
const shutdownTimeout = 10 * time.Second

// Run serves until the listener fails or SIGINT/SIGTERM is received, then
// drains in-flight requests so the caller can close the database cleanly.
func Run(port string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return serve(ctx, &http.Server{
		Addr:              ":" + port,
		Handler:           NewRouter(),
		ReadHeaderTimeout: 10 * time.Second,
	})
}

func serve(ctx context.Context, srv *http.Server) error {
	log.Printf("The server will be running on %s", srv.Addr)
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}
	log.Println("Shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
