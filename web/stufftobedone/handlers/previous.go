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
)

func PreviousHandler(c *gin.Context) {
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

  c.HTML(http.StatusOK, "previous.html", gin.H{
        "title": "Previous to " + dayAsString,
        "dayAsString": dayAsString,
        "todayAsString": todayAsString,
        "tomorrowAsString": tomorrow.DateAsString,
        "yesterdayAsString": yesterday.DateAsString,
        "logouturl" : user.LogoutUrl,
        "bookID": bookID,
        "isProduction": user.IsProduction,
  })
}



func ApiPreviousHandler(c *gin.Context) {
  result := []domain.OpenDay{}

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

  days, err := dayRepository.OpenItemsBeforeDay(bookID, dayAsInt)
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
    dayTask.Build(days[i], tasks[i])
    // find the dateAsInt in the result
    found := false
    for index, openDay := range result {
      if openDay.DateAsInt == dayTask.DateAsInt {
        result[index].DayTasks = append(openDay.DayTasks, dayTask)
        found = true
        log.Println("add to existing", dayTask.DateAsInt)
      }
	  }
    if !found {
      openDay := domain.OpenDay{ DateAsInt: dayTask.DateAsInt, Display: common.ConvertIntToDisplay(dayTask.DateAsInt), DayTasks: []domain.DayTask{ dayTask } }
      result = append(result, openDay)
        log.Println("add new", dayTask.DateAsInt)
    }
  }

  if err != nil {
    JsonError(c, err)
  } else {
    c.JSON(http.StatusOK, gin.H{
        "data": result,
    })
  }
} 