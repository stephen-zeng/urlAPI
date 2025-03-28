package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"urlAPI/request"
)

func sessionHandler(c *gin.Context) {
	var sessionRequest request.Request
	if err := sessionBuilder(c, &sessionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := sessionRequest.Processor.Session.Process(&sessionRequest.DB.Session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, sessionRequest.Processor.Session)
	}
}

func sessionBuilder(c *gin.Context, r *request.Request) error {
	c.Header("Access-Control-Allow-Origin", "*")
	err := c.ShouldBind(&r.Processor.Session)
	if err != nil { // auth err
		return err
	}
	r.DB.Session.Token = c.Request.Header.Get("Authorization")
	r.DB.Session.Term = r.Processor.Session.LoginTerm
	r.Processor.Session.SessionIP = c.ClientIP()
	return nil
}
