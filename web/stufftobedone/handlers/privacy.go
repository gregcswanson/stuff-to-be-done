package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func PrivacyHandler(c *gin.Context) {
  user := GetAppUser(c)
  c.HTML(http.StatusOK, "privacy.html", gin.H{
        "title": "Privacy Policy",
        "logouturl" : user.LogoutUrl,
        "isLoggedIn": user.IsLoggedIn,
  })
}
