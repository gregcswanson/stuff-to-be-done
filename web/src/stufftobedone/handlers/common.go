package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stufftobedone/domain"
)

func TestErrorPage(c *gin.Context) {
	err := errors.New("Test error")
	ErrorPage(c, err)
}

func ErrorPage(c *gin.Context, err error) {
	// return the result
	c.HTML(http.StatusOK, "error.html", gin.H{
		"title":   "Oh crap, something went wrong",
		"message": err.Error(),
	})
}

func JsonError(c *gin.Context, err error) {
	log.Println("JsonError: " + err.Error())
	c.JSON(500, gin.H{
		"message": err.Error(),
	})
}

func JsonErrorMessage(c *gin.Context, err string) {
	c.JSON(500, gin.H{
		"message": err,
	})
}

// helper function for the handers
func GetAppUser(c *gin.Context) domain.AppUser {
	appUser := c.MustGet("user").(domain.AppUser)
	return appUser
}
