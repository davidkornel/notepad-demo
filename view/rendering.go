package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func Render(c *gin.Context, data gin.H, templateName string) {
	// Authentication should be checked here

	acceptHeader := c.Request.Header.Get("Accept")

	switch {
	case strings.Contains(acceptHeader, "application/json"):
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}
