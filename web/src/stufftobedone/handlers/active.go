package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"stufftobedone/common"
	"stufftobedone/domain"
	"stufftobedone/repositories"
	//"time"
)

func ApiActiveHandler(c *gin.Context) {
	result := []domain.OpenDay{}

	// get all for today, an incomplete for earlier
	// return as a grouped by  day result
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
	// get today
	days, err := dayRepository.Day(bookID, dayAsInt)
	if err != nil {
		JsonError(c, err)
		return
	}
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
		if !tasks[i].IsDeleted {
			for index, openDay := range result {
				if openDay.DateAsInt == dayTask.DateAsInt {
					result[index].DayTasks = append(openDay.DayTasks, dayTask)
					found = true
					log.Println("add to existing", dayTask.DateAsInt)
				}
			}
			if !found {
				openDay := domain.OpenDay{DateAsInt: dayTask.DateAsInt, Display: common.ConvertIntToDisplay(dayTask.DateAsInt), DayTasks: []domain.DayTask{dayTask}}
				result = append(result, openDay)
			}
		}
	}
	// if there are no current active items return an empty day
	if len(tasks) == 0 {
		openDay := domain.OpenDay{DateAsInt: dayAsInt, Display: common.ConvertIntToDisplay(dayAsInt), DayTasks: []domain.DayTask{}}
		result = append(result, openDay)
	}

	// get the history
	days, err = dayRepository.OpenItemsBeforeDay(bookID, dayAsInt)
	if err != nil {
		JsonError(c, err)
		return
	}
	// get the tasks
	taskIds = []string{}
	for i := 0; i < len(days); i++ {
		taskIds = append(taskIds, days[i].TaskID)
	}
	tasks, err = taskRepository.FindByIds(bookID, taskIds)
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
			openDay := domain.OpenDay{DateAsInt: dayTask.DateAsInt, Display: common.ConvertIntToDisplay(dayTask.DateAsInt), DayTasks: []domain.DayTask{dayTask}}
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
