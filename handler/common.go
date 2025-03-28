package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"urlAPI/request"
)

func taskSaver(r *request.Request) {
	if r.Security.General.SkipDB {
		return
	}
	if err := r.DB.Task.Create(); err != nil {
		log.Println(err)
	}
}

func returner(c *gin.Context, jsonString, url string) {
	if c.Query("format") == "json" {
		c.JSON(http.StatusOK, jsonString)
	} else {
		c.Redirect(http.StatusFound, url)
	}
}
