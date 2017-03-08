package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func FourZeroFourHandler(c *gin.Context) {
  c.HTML(http.StatusOK, "404.html", gin.H{
        "title": "404",
  })
}

func FourZeroFourApiHandler(c *gin.Context) {
  c.JSON(404, gin.H{
      "message": "unknown endpoint",
  })
}
