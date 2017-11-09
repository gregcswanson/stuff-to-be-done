package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
