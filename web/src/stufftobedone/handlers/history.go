package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HistoryHandler(c *gin.Context) {
	user := GetAppUser(c)
	bookID := c.Param("book")

	c.HTML(http.StatusOK, "history.html", gin.H{
		"title":        "History",
		"logouturl":    user.LogoutUrl,
		"isLoggedIn":   true,
		"bookID":       bookID,
		"isProduction": user.IsProduction,
	})
}
