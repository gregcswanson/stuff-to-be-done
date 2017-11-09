package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AboutHandler(c *gin.Context) {
	user := GetAppUser(c)
	c.HTML(http.StatusOK, "about.html", gin.H{
		"title":        "About",
		"logouturl":    user.LogoutUrl,
		"isLoggedIn":   user.IsLoggedIn,
		"isProduction": user.IsProduction,
	})
}
