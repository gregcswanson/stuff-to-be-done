package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "stufftobedone/repositories"
  "time"
)

func CompleteHandler(c *gin.Context) {
  user := GetAppUser(c)
  bookID := c.Param("book")
  dayAsString := c.Param("day")
  c.HTML(http.StatusOK, "complete.html", gin.H{
        "title": "No doing",
        "logouturl" : user.LogoutUrl,
        "dayAsString": dayAsString,
        "bookID": bookID,
  })
}

func ApiCompleteHandler(c *gin.Context) {
  // get a list of copmlete elements for this book
  bookID := c.Param("bookId")
  r := repositories.NewTaskRepository(c.Request)
  laterItems, err := r.Complete(bookID)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
  } else {
    c.JSON(http.StatusOK, gin.H{
        "data": laterItems,
    })
  }
}

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

