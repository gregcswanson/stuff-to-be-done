package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"stufftobedone/repositories"
	"stufftobedone/usecases"
)

func TrashHandler(c *gin.Context) {
	user := GetAppUser(c)
	bookID := c.Param("book")
	c.HTML(http.StatusOK, "trash.html", gin.H{
		"title":        "Trash",
		"logouturl":    user.LogoutUrl,
		"bookID":       bookID,
		"isProduction": user.IsProduction,
	})
}

func ApiTrashGetHandler(c *gin.Context) {
	taskUseCases := usecases.NewTaskUseCases(c)
	bookID := c.Param("bookId")

	groupedDayTasks, err := taskUseCases.FindTrashGroupedByDay(bookID)
	// DEVELOPER CODE - TESTING ERRORS
	//JsonErrorMessage(c, "testing error")
	// ~ DEVELOPER CODE
	if err != nil {
		JsonError(c, err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": groupedDayTasks,
		})
	}
}

func ApiTrashEmptyHandler(c *gin.Context) {
	taskUseCases := usecases.NewTaskUseCases(c)

	bookID := c.Param("bookId")

	err := taskUseCases.EmptyTrash(bookID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": "OK",
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
