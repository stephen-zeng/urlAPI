package router

import (
	"github.com/gin-gonic/gin"
)

var r *gin.Engine = gin.Default()

func Start() {
	setAPI()
	r.Run(":8080")
}
