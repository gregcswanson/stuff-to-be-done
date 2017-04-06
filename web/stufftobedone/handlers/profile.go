package handlers

import (
  "net/http"
  //"log"
  "github.com/gin-gonic/gin"
  "stufftobedone/repositories"
  "time"
)

func ProfileIndexsHandler(c *gin.Context) {
  // setup the required repositories
  user := GetAppUser(c)
  bookRepository := repositories.NewBookRepository(c.Request)
  // get or create the default book for this user
  defaultBook, err := bookRepository.GetDefault(user)
  if err != nil {
    c.JSON(500, gin.H{
          "message": err,
        })
  }
  // redirect to the current day
  todayAsString := time.Now().Format("20060102")
  c.Redirect(http.StatusFound, "/book/" + defaultBook.ID + "/day/" + todayAsString)
  //c.JSON(http.StatusOK, gin.H{
  //  "bookID": defaultBook.ID,
  //})
}