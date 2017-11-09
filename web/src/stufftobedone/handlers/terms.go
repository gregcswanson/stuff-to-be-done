package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TermsHandler(c *gin.Context) {
	user := GetAppUser(c)
	c.HTML(http.StatusOK, "terms.html", gin.H{
		"title":        "Terms",
		"logouturl":    user.LogoutUrl,
		"isLoggedIn":   user.IsLoggedIn,
		"isProduction": user.IsProduction,
	})
}
