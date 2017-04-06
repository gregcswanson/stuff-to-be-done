package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "stufftobedone/repositories"
  "stufftobedone/domain"
  "time"
  "log"
  "strconv"
)

func LaterHandler(c *gin.Context) {
  user := GetAppUser(c)
  bookID := c.Param("book")
  //dayAsString := c.Param("day")
  c.HTML(http.StatusOK, "later.html", gin.H{
        "title": "Do later",
        "logouturl" : user.LogoutUrl,
        //"dayAsString": dayAsString,
        "bookID": bookID,
  })
}

func BookLaterHandler(c *gin.Context) {
  // get a list of available elements for this book
  bookID := c.Param("bookId")
  log.Println("getting list of later items for the book", bookID)
  r := repositories.NewTaskRepository(c.Request)
  laterItems, err := r.Later(bookID)
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

func BookLaterPostHandler(c *gin.Context) {
  // setup the required repositories
  user := GetAppUser(c)
  // bookRepository := repositories.NewBookRepository(c.Request)
  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  bookID := c.Param("bookId")
  elementName := c.PostForm("elementName")

  // validate element and book

  task, err := taskRepository.Create(elementName, bookID, 0, "", user.ID)
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

func BookLaterPutHandler(c *gin.Context) {
  // setup the required repositories
  user := GetAppUser(c)
  // bookRepository := repositories.NewBookRepository(c.Request)
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
  log.Println("Original task", task)

  // validate bookId
  if task.BookID != bookID {
    c.JSON(500, gin.H{
      "message": "BookID is not valid",
    })
    return
  }
  
  if !task.LaterCanUpdate() {
    c.JSON(500, gin.H{
      "message": "Task cannot be updated",
    })
    return
  }

  task.Data = c.PostForm("Data")
  if (!task.IsCompleted && c.PostForm("IsCompleted") == "true") {
    task.Completed = time.Now()
  }
  task.IsCompleted = c.PostForm("IsCompleted") == "true"
  
  task.Updated = time.Now()
  task.UpdatedBy = user.ID
  log.Println("Updated task", task)

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

func BookLaterDeleteHandler(c *gin.Context) {
  // setup the required repositories
  
  log.Println("delete request handler")

  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  user := GetAppUser(c)
  bookID := c.Param("bookId")
  taskID := c.Param("taskId")
  log.Println("delete", bookID, taskID)
  
  // validate element and book
  task, err := taskRepository.FindById(bookID, taskID)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  }
  
  if !task.LaterCanDelete() {
    c.JSON(500, gin.H{
      "message": "Task cannot be deleted",
    })
    return
  }

  task.IsDeleted = true
  task.Updated = time.Now()
  task.UpdatedBy = user.ID
  log.Println("Send the task to trash", task)

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

func BookLaterRestoreHandler(c *gin.Context) {
  // setup the required repositories
  
  log.Println("undo delete request handler")

  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  user := GetAppUser(c)
  bookID := c.Param("bookId")
  taskID := c.Param("taskId")
  log.Println("undo delete", bookID, taskID)
  
  // validate element and book
  task, err := taskRepository.FindById(bookID, taskID)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  }
  
  if !task.LaterCanUndoDelete() {
    c.JSON(500, gin.H{
      "message": "Task delete cannot be undone",
    })
    return
  }

  task.IsDeleted = false
  task.Updated = time.Now()
  task.UpdatedBy = user.ID
  log.Println("Restore the task from trash", task)

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

func BookLaterDoTodayHandler(c *gin.Context) {
 
  // setup the required repositories
  user := GetAppUser(c)
  dayRepository := repositories.NewDayRepository(c.Request)
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
  log.Println("Original task", task)

  // validate bookId
  if task.BookID != bookID {
    c.JSON(500, gin.H{
      "message": "BookID is not valid",
    })
    return
  }
  
  if !task.LaterCanUpdate() {
    c.JSON(500, gin.H{
      "message": "Task cannot be updated",
    })
    return
  }

  task.Data = c.PostForm("Data")
  
  task.Updated = time.Now()
  task.UpdatedBy = user.ID
  log.Println("Updated task", task)

  // add it to today
  todayAsString := time.Now().Format("20060102")
  dayAsInt, _ := strconv.Atoi(todayAsString)

  // create or find the day record
  day, err := dayRepository.FindByTaskId(bookID, task.ID,  dayAsInt)
  if err != nil {
    // not task was found on this day, create it
    day := domain.Day{BookID: bookID, TaskID: task.ID, Data: task.Data, DateAsInt: dayAsInt, IsActioned: false, IsCompleted: false, Created: time.Now()}
    day, err = dayRepository.Create(day)
    if err != nil {
      c.JSON(500, gin.H{
        "message": err,
      })
      return
    }
  } else {
    // update the day task
    day.Data = task.Data
    day.IsCompleted = task.IsCompleted
    day.IsActioned = false
    day, err = dayRepository.Update(day)
    if err != nil {
      c.JSON(500, gin.H{
        "message": err,
      })
      return
    }
  }

  // update the task
  task.CurrentDate = dayAsInt
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

