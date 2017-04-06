package usecases

import (
	//"net/http"
	"github.com/gin-gonic/gin"
	"stufftobedone/domain"
	//"appengine"
	//"appengine/datastore"
	"stufftobedone/repositories"
  "stufftobedone/common"
	"errors"
	"time"
	"log"
)

type TaskUseCases struct {
	context    *gin.Context
	DayRepository *repositories.DayRepository
	TaskRepository *repositories.TaskRepository
	User domain.AppUser
}

func NewTaskUseCases(c *gin.Context) *TaskUseCases {
	r := new(TaskUseCases)
	r.context = c
	r.User = c.MustGet("user").(domain.AppUser)
	r.DayRepository = repositories.NewDayRepository(r.context.Request)
  	r.TaskRepository = repositories.NewTaskRepository(r.context.Request)
	return r
}

func (r *TaskUseCases) DoLater(bookID string, taskID string, dayID string, data string) (domain.DayTask, error) {
	result := domain.DayTask{}

	task, err := r.TaskRepository.FindById(bookID, taskID)
	if err != nil {
		log.Println("task not found")
		return result, err
	}
  
	// validate bookId
	if task.BookID != bookID {
		log.Println("book not found")
		return result, errors.New("BookID is not valid")
	}

	// if the task is completed it is not on a day, so send it back to do later
	if task.IsCompleted {
		task.IsCompleted = false
  	task.Updated = time.Now()
  	task.UpdatedBy = r.User.ID

  	task, err = r.TaskRepository.Update(task)
		if err != nil {
			return result, err
		}
		result.BuildTask(task)
		return result, nil
	}

	// get the day item
	var day domain.Day
	if (dayID == ""){
  	day, err = r.DayRepository.FindByTaskId(bookID, task.ID,  task.CurrentDate)
	} else {
		day, err = r.DayRepository.FindById(bookID, dayID)
	}
	if err != nil {
		log.Println("day not found")
		return result, errors.New("Day item not found")
	}
	log.Println("Original day", day)
  
	day.Data = data
	task.Data = data
	day.IsActioned = true
	task.CurrentDate = 0
	task.Updated = time.Now()
	task.UpdatedBy = r.User.ID
	
		log.Println("save task")
	task, err = r.TaskRepository.Update(task)
	if err != nil {
		return result, err
	} 
		log.Println("save day")
  day, err = r.DayRepository.Update(day)
	result.Build(day, task)
  return result, err
}

func (r *TaskUseCases) DoOnDate(bookID string, taskId string, dateAsString string, data string) (domain.Task, domain.Day, error) {
	day := domain.Day{}

  // validate element and book
  task, err := r.TaskRepository.FindById(bookID, taskId)
  if err != nil || task.BookID != bookID {
    return task, day, errors.New("BookID is not valid")
  }
  
  if !task.LaterCanUpdate() {
    return task, day, errors.New("Task cannot be updated")
  }

  convertedDate, err := common.ConvertStringToDates(dateAsString)
  if err != nil {
    return task, day, errors.New("Date " + dateAsString + " is not valid")
  }
	log.Println(convertedDate)

  task.Data = data
  task.Updated = time.Now()
  task.UpdatedBy = r.User.ID
	task.IsDeleted = false


  // create or find the day record
  day, err = r.DayRepository.FindByTaskId(bookID, task.ID,  convertedDate.DateAsInt)
  if err != nil {
    // not task was found on this day, create it
    day := domain.Day{BookID: bookID, TaskID: task.ID, Data: task.Data, DateAsInt: convertedDate.DateAsInt, IsActioned: false, IsCompleted: false, Created: time.Now()}
    day, err = r.DayRepository.Create(day)
    if err != nil {
      return task, day, err
    }
  } else {
    // update the day task
    day.Data = task.Data
    day.IsCompleted = task.IsCompleted
    day.IsActioned = false
    day, err = r.DayRepository.Update(day)
    if err != nil {
      return task, day, err
    }
  }

  // update the task
  task.CurrentDate = convertedDate.DateAsInt
  task, err = r.TaskRepository.Update(task)
  return task, day, err
}




