package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func HistoryHandler(c *gin.Context) {
  user := GetAppUser(c)
  bookID := c.Param("book")
  
  c.HTML(http.StatusOK, "history.html", gin.H{
        "title": "History",
        "logouturl" : user.LogoutUrl,
        "isLoggedIn": true,
        "bookID": bookID,
  })
}

