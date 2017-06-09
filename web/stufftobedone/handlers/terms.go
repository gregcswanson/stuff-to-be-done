package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func TermsHandler(c *gin.Context) {
  user := GetAppUser(c)
  c.HTML(http.StatusOK, "terms.html", gin.H{
        "title": "Terms",
        "logouturl" : user.LogoutUrl,
        "isLoggedIn": user.IsLoggedIn,
  })
}
