package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "stufftobedone/repositories"
  //"time"
)

func BookIndexsHandler(c *gin.Context) {
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
  //todayAsString := time.Now().Format("20060102")
  c.Redirect(http.StatusFound, "/book/" + defaultBook.ID + "/live") //"/"  + todayAsString + "/day")
}

func BookLiveHander(c *gin.Context) {
  user := GetAppUser(c)
  bookID := c.Param("book")
  
  c.HTML(http.StatusOK, "book.html", gin.H{
        "title": "Current",
        "logouturl" : user.LogoutUrl,
        "bookID": bookID,
  })
}
