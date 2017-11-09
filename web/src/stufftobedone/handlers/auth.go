package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/profile")
}

func LogoutHandler(c *gin.Context) {
	user := GetAppUser(c)
	// TO DO - check this status is correct
	if user.IsLoggedIn {
		c.Redirect(http.StatusMovedPermanently, user.LogoutUrl)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}
