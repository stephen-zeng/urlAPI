package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"urlAPI/request"
)

type JsonStruct struct {
	Prompt         string `json:"prompt"`
	OriginalPrompt string `json:"original_prompt"`
	ActualPrompt   string `json:"actual_prompt"`
	Response       string `json:"response"`
	URL            string `json:"url"`
}

func taskSaver(r *request.Request) {
	if r.Security.General.SkipDB {
		return
	}
	if err := r.DB.Task.Create(); err != nil {
		log.Println(err)
	}
}

func returner(c *gin.Context, jsonString, url string) {
	var jsonStruct JsonStruct
	json.Unmarshal([]byte(jsonString), &jsonStruct)
	if c.Query("format") == "json" {
		c.JSON(http.StatusOK, jsonStruct)
	} else {
		c.Redirect(http.StatusFound, url)
	}
}

func getScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return `https://`
	}
	if scheme := c.GetHeader("X-Forwarded-Proto"); scheme != "" {
		return scheme + `://`
	}
	return `http://`
}
