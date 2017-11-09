package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PrivacyHandler(c *gin.Context) {
	user := GetAppUser(c)
	c.HTML(http.StatusOK, "privacypolicy.htm", gin.H{
		"title":        "Privacy Policy",
		"logouturl":    user.LogoutUrl,
		"isLoggedIn":   user.IsLoggedIn,
		"isProduction": user.IsProduction,
	})
}
