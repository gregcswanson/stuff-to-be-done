package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "stufftobedone/repositories"
  "stufftobedone/domain"
  "time"
  "strconv"
  "log"
  "stufftobedone/common"
  "stufftobedone/usecases"
)

func DayHandler(c *gin.Context) {
  user := GetAppUser(c)
  bookID := c.Param("book")
  dayAsString := c.Param("day")
  // Convert this to a date, then work out yesterday and tomorrow
  convertedDate, err := common.ConvertStringToDates(dayAsString)
  if err != nil {
    ErrorPage(c, err)
    return
  }
  tomorrow, _ := convertedDate.Tomorrow()
  yesterday, _ := convertedDate.Yesterday()
  todayAsString := time.Now().Format("20060102")

  c.HTML(http.StatusOK, "day.html", gin.H{
        "title": dayAsString,
        "dayAsString": dayAsString,
        "todayAsString": todayAsString,
        "tomorrowAsString": tomorrow.DateAsString,
        "yesterdayAsString": yesterday.DateAsString,
        "logouturl" : user.LogoutUrl,
        "bookID": bookID,
  })
}

func ApiDayNewElementHandler(c *gin.Context) {
  taskUseCases := usecases.NewTaskUseCases(c);
  
  bookID := c.Param("bookId")
  dayAsString := c.Param("dayAsString")
  elementName := c.PostForm("elementName")

  dayTask, err := taskUseCases.NewTaskOnDate(bookID, dayAsString, elementName)
  if err != nil {
    JsonError(c, err)
    return
  } else {
    c.JSON(http.StatusOK, gin.H{
        "data": dayTask,
    })
  }
  
/*
  // setup the required repositories
  user := GetAppUser(c)
  dayRepository := repositories.NewDayRepository(c.Request)
  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  bookID := c.Param("bookId")
  dayAsString := c.Param("dayAsString")
  dayAsInt, err := strconv.Atoi(dayAsString)
  if err != nil {
    JsonErrorMessage(c, "Date is invalid")
    return
  }
  elementName := c.PostForm("elementName")

  // create the task, and add it to the requested date
  task, err := taskRepository.Create(elementName, bookID, dayAsInt, "", user.ID)
  if err != nil {
    JsonError(c, err)
    return
  } 
  // create the linked day record
  day := domain.Day{BookID: bookID, TaskID: task.ID, Data: task.Data, DateAsInt: dayAsInt, IsActioned: false, IsCompleted: false, Created: time.Now()}
  day, err = dayRepository.Create(day)

  // build the result record
  dayTask := domain.DayTask{}
  dayTask.Build(day, task)

  // return the new day record
  if err != nil {
    JsonError(c, err)
  } else {
    c.JSON(http.StatusOK, gin.H{
        "data": dayTask,
    })
  }
*/
}

func ApiDayHander(c *gin.Context) {
  result := []domain.DayTask{}

  //user := GetAppUser(c)
  dayRepository := repositories.NewDayRepository(c.Request)
  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  bookID := c.Param("bookId")
  dayAsString := c.Param("dayAsString")
  dayAsInt, err := strconv.Atoi(dayAsString)
  if err != nil {
    JsonErrorMessage(c, "Date is invalid")
    return
  }

  days, err := dayRepository.Day(bookID, dayAsInt)
  if err != nil {
    JsonError(c, err)
    return
  }
  // get the tasks
  taskIds := []string{}
  for i := 0; i < len(days); i++ {
    taskIds = append(taskIds, days[i].TaskID)
  }
  tasks, err := taskRepository.FindByIds(bookID, taskIds)
  if err != nil {
    JsonError(c, err)
    return
  }
  for i := 0; i < len(tasks); i++ {
    dayTask := domain.DayTask{}
    log.Println(days[i], tasks[i])
    dayTask.Build(days[i], tasks[i])
    log.Println("daytask", dayTask)
    result = append(result, dayTask)
  }

  if err != nil {
    JsonError(c, err)
  } else {
    c.JSON(http.StatusOK, gin.H{
        "data": result,
    })
  }
}

func ApiDayPutHandler(c *gin.Context) {
  
  // setup the required repositories
  user := GetAppUser(c)
  dayRepository := repositories.NewDayRepository(c.Request)
  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  bookID := c.Param("bookId")
  
  // validate day and book
  task, err := taskRepository.FindById(bookID, c.PostForm("TaskID"))
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

  // get the day item
  day, err := dayRepository.FindById(bookID, c.PostForm("DayID"))
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  }
  log.Println("Original day", day)
  
  day.Data = c.PostForm("Data")
  task.Data = c.PostForm("Data")
  if (task.CurrentDate == day.DateAsInt){
    if (!task.IsCompleted && c.PostForm("IsCompleted") == "true") {
      task.Completed = time.Now()
    } else if(task.IsCompleted && c.PostForm("IsCompleted") == "false") {
      task.Completed = time.Now()
    }
    task.IsCompleted = c.PostForm("IsCompleted") == "true"
    day.IsCompleted = task.IsCompleted
  }
  day.Tags = c.PostForm("Tags")
  task.Tags = day.Tags
  
  task.Updated = time.Now()
  task.UpdatedBy = user.ID
  // if this day is the current day for the task, toggle the task coomplete
  

  task, err = taskRepository.Update(task)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  } 
  day, err = dayRepository.Update(day)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
  } else {
    // build the result record
    dayTask := domain.DayTask{}
    dayTask.Build(day, task)

    c.JSON(http.StatusOK, gin.H{
        "data": dayTask,
    })
  }

}

func ApiDayDeleteHandler(c *gin.Context) {
  // setup the required repositories
  user := GetAppUser(c)
  dayRepository := repositories.NewDayRepository(c.Request)
  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  bookID := c.Param("bookId")
  taskID := c.Param("taskId")
  //dayAsString := c.Param("dayAsString")
  //data := c.PostForm("Data")

  //convertedDate, _ := common.ConvertStringToDates(dayAsString)
  
  // validate day and book
  task, err := taskRepository.FindById(bookID, taskID)
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

  // get the day item is the task is currently on a day
  day := domain.Day{ BookID: task.BookID, DateAsInt: 0, Data: task.Data, IsCompleted: task.IsCompleted, Created: time.Now() }
  if task.CurrentDate != 0 {
    day, err := dayRepository.FindByTaskId(bookID, taskID, task.CurrentDate)
    if err != nil {
      c.JSON(500, gin.H{
        "message": err,
      })
      return
    }
    day.IsActioned = true
    //day.Data = data
    task.CurrentDate = 0
	  task.LastDayID = day.ID

    day, err = dayRepository.Update(day)
    if err != nil {
      c.JSON(500, gin.H{
        "message": err,
      })
      return
    } 
  }

  //task.Data = data
  task.Updated = time.Now()
  task.UpdatedBy = user.ID
  task.IsDeleted = true

  task, err = taskRepository.Update(task)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  } 

  // build the result record
  dayTask := domain.DayTask{}
  dayTask.Build(day, task)

  c.JSON(http.StatusOK, gin.H{
      "data": dayTask,
  })
}

func ApiDayDoLaterPutHandler(c *gin.Context) {
  taskUseCases := usecases.NewTaskUseCases(c);
  
  bookID := c.Param("bookId")
  taskID := c.PostForm("TaskID")
  dayID := c.PostForm("DayID")
  data :=  c.PostForm("Data")

  result, err := taskUseCases.DoLater(bookID, taskID, dayID, data)
  if err != nil {
    JsonError(c, err)
  } else {
    c.JSON(http.StatusOK, gin.H{
        "data": result,
    })
  }

}

func CommentPutHandler(c *gin.Context) {
  // setup the required repositories
  log.Println("put comment");
  //user := GetAppUser(c)
  dayRepository := repositories.NewDayRepository(c.Request)
  taskRepository := repositories.NewTaskRepository(c.Request)

  // get the parameters
  bookID := c.Param("bookId")
  taskID := c.PostForm("TaskID")
  dayID := c.PostForm("DayID")
  comment :=  c.PostForm("Comment")

  // get the task
  task, err := taskRepository.FindById(bookID, taskID)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  }
  // get the day item
  day, err := dayRepository.FindById(bookID, dayID)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  }
  
  day.Comment = comment
  day, err = dayRepository.Update(day)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
  } else {
    // build the result record
    dayTask := domain.DayTask{}
    dayTask.Build(day, task)
    log.Println(dayTask)
    c.JSON(http.StatusOK, gin.H{
        "data": dayTask,
    })
  }
}

func ApiDayDoOnDatePutHander(c *gin.Context) {
  taskUseCases := usecases.NewTaskUseCases(c)

  // get the parameters
  bookID := c.Param("bookId")
  taskID := c.PostForm("TaskID")
  dateAsString := c.PostForm("DayAsString")
  data := c.PostForm("Data")

  task, day, err := taskUseCases.DoOnDate(bookID, taskID, dateAsString, data)
  if err != nil {
    c.JSON(500, gin.H{
      "message": err,
    })
    return
  }
  // build the result record
  dayTask := domain.DayTask{}
  dayTask.Build(day, task)

  c.JSON(http.StatusOK, gin.H{
      "data": dayTask,
  })
}

