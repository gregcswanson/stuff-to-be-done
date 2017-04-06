package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func AboutHandler(c *gin.Context) {
  user := GetAppUser(c)
  c.HTML(http.StatusOK, "about.html", gin.H{
        "title": user.Name,
        "logouturl" : user.LogoutUrl,
  })
}
