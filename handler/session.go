package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"urlAPI/request"
	"urlAPI/util"
)

func sessionHandler(c *gin.Context) {
	var sessionRequest request.Request
	if err := sessionBuilder(c, &sessionRequest); err != nil {
		util.ErrorPrinter(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := sessionRequest.Processor.Session.Process(&sessionRequest.DB.Session); err != nil {
		log.Printf("%s from %s\n", err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, sessionRequest.Processor.Session)
	}
}

func sessionBuilder(c *gin.Context, r *request.Request) error {
	c.Header("Access-Control-Allow-Origin", "*")
	if err := c.ShouldBind(&r.Processor.Session); err != nil { // auth Error
		return err
	}
	r.DB.Session.Token = c.Request.Header.Get("Authorization")
	r.DB.Session.Term = r.Processor.Session.LoginTerm
	r.Processor.Session.SessionIP = c.ClientIP()
	return nil
}
