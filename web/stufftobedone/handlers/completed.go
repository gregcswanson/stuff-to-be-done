package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "stufftobedone/repositories"
 // "time"
)

func CompletedHandler(c *gin.Context) {
  user := GetAppUser(c)
  bookID := c.Param("book")
  dayAsString := c.Param("day")
  c.HTML(http.StatusOK, "completed.html", gin.H{
        "title": "No doing",
        "logouturl" : user.LogoutUrl,
        "dayAsString": dayAsString,
        "bookID": bookID,
  })
}

func ApiCompletedHandler(c *gin.Context) {
  // get a list of copmlete elements for this book
  bookID := c.Param("bookId")
  r := repositories.NewTaskRepository(c.Request)
  items, err := r.Completed(bookID)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
  } else {
    c.JSON(http.StatusOK, gin.H{
        "data": items,
    })
  }
}

/*
func ApiFromCompleteToDoLaterHandler(c *gin.Context) {
  // setup the required repositories
  user := GetAppUser(c)
  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  bookID := c.Param("bookId")
  
  // validate element and book
  task, err := taskRepository.FindById(bookID, c.PostForm("ID"))
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  }
  
  if !task.LaterCanUpdate() {
    c.JSON(500, gin.H{
      "message": "Task cannot be updated",
    })
    return
  }

  task.IsComplete = false
  task.Updated = time.Now()
  task.UpdatedBy = user.ID

  task, err = taskRepository.Update(task)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
  } else {
    c.JSON(http.StatusOK, gin.H{
        "data": task,
    })
  }

}
*/

