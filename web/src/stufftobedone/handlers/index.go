package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexHandler(c *gin.Context) {
	user := GetAppUser(c)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":            "The Founder",
		"showregistration": "true",
		"isLoggedIn":       user.IsLoggedIn,
		"logouturl":        user.LogoutUrl,
		"isProduction":     user.IsProduction,
	})
}
