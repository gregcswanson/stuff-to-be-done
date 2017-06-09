package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
      user := GetAppUser(c)
      c.HTML(http.StatusOK, "index.html", gin.H{
        "title": "The Founder",
        "showregistration": "true",
        "isLoggedIn": user.IsLoggedIn,
        "logouturl" : user.LogoutUrl,
      })
}
