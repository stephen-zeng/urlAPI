package main

import (
	"backend/internal/data"
	"backend/internal/router"
)

func main() {
	data.InitSetting(data.DataConfig())
	router.HttpServer()
}
