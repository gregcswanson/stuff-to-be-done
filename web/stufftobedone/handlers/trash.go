package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "stufftobedone/repositories"
)

func TrashHandler(c *gin.Context) {
  user := GetAppUser(c)
  bookID := c.Param("book")
  c.HTML(http.StatusOK, "trash.html", gin.H{
        "title": "Trash",
        "logouturl" : user.LogoutUrl,
        "bookID": bookID,
  })
}

func ApiTrashGetHandler(c *gin.Context) {
  // get a list of copmlete elements for this book
  bookID := c.Param("bookId")
  r := repositories.NewTaskRepository(c.Request)
  items, err := r.Deleted(bookID)
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

