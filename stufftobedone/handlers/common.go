package handlers

import (
	"net/http"
  "github.com/gin-gonic/gin"
  "stufftobedone/domain"
  "errors"
)

func TestErrorPage(c *gin.Context) {
  err := errors.New("Test error")
  ErrorPage(c, err)
}

func ErrorPage(c *gin.Context, err error) {
  // return the result
  c.HTML(http.StatusOK, "error.html", gin.H{
    "title": "Oh crap, something went wrong",
    "message": err,
  })
}

func JsonError(c *gin.Context, err error) {
  c.JSON(500, gin.H{
      "message": err,
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
