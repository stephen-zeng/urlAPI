package router

import (
	"github.com/gin-gonic/gin"
)

var r *gin.Engine = gin.Default()

func Start() {
	setAPI()
	//sessionListener()
	r.Run(":8080")
}
