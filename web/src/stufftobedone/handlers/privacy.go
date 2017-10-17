package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func PrivacyHandler(c *gin.Context) {
  user := GetAppUser(c)
  c.HTML(http.StatusOK, "privacypolicy.htm", gin.H{
        "title": "Privacy Policy",
        "logouturl" : user.LogoutUrl,
        "isLoggedIn": user.IsLoggedIn,
        "isProduction": user.IsProduction,
  })
}
