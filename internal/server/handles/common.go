package handles

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"urlAPI/internal/op"
	"urlAPI/internal/server/trust"
)

// returner sends a generation result. Results carry server-relative URLs
// so that cached results never embed another request's Host header; they
// are made absolute here for JSON consumers on other origins.
func returner(c *gin.Context, result op.GenerateResult) {
	if c.Query("format") == "json" {
		result.URL = trust.Get().Absolute(c.Request, result.URL)
		c.JSON(http.StatusOK, result)
	} else {
		c.Redirect(http.StatusFound, result.URL)
	}
}

func errorReturner(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
