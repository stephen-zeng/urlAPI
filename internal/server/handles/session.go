package handles

import (
	"log"
	"net/http"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/internal/op"
	"urlAPI/util"

	"github.com/gin-gonic/gin"
)

func SessionHandler(c *gin.Context) {
	allowedIPs := database.SettingsStore.Get().Security.DashboardAllowedIPs
	clientIP := c.ClientIP()
	// An empty list imposes no restriction (the default is "*").
	if len(allowedIPs) > 0 && !util.WildcardChecker(&allowedIPs, &clientIP) {
		log.Printf("Dashboard access denied for %s\n", clientIP)
		c.JSON(http.StatusForbidden, gin.H{"error": "IP " + clientIP + " is not allowed to access the dashboard"})
		return
	}
	var request op.Session
	if err := c.ShouldBind(&request); err != nil { // auth Error
		util.ErrorPrinter(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authSession := model.Session{
		Token: c.Request.Header.Get("Authorization"),
		Term:  request.LoginTerm,
	}
	request.SessionIP = c.ClientIP()
	request.SessionToken = c.Request.Header.Get("Authorization")

	response, err := op.HandleSession(request, authSession)
	if err != nil {
		log.Printf("%s from %s\n", err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
