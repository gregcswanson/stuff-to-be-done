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
	"strconv"
	"sort"
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

func (r *TaskUseCases) FindTasksSetAsLaterGroupedByDay(bookID string) ([]domain.GroupedDayTask, error) {
	result := []domain.GroupedDayTask{}

	// get a list of all the open tasks
	tasks, err := r.TaskRepository.Later(bookID)
	if err != nil {
    return result, err
  }

	// get the last open day from the last day ID for each task
	dayIds := []string{}
  for i := 0; i < len(tasks); i++ {
    dayIds = append(dayIds, tasks[i].LastDayID)
  }
	days, err := r.DayRepository.FindByIds(bookID, dayIds)
	if err != nil {
    return result, err
  }

	// build a list of day tasks
	dayTasks := []domain.DayTask{}
	for i := 0; i < len(tasks); i++ {
    dayTask := domain.DayTask{}
    dayTask.Build(days[i], tasks[i])
		dayTasks = append(dayTasks, dayTask)
  }

	// sort the list by date descending
	sort.Sort(domain.DayTaskSort(dayTasks))

	// create a list of laterDays and add group the day tasks by unquie days
  for i := 0; i < len(dayTasks); i++ {
    // find the dateAsInt in the result
    found := false
    for index, groupedDayTask := range result {
      if groupedDayTask.DateAsInt == dayTasks[i].DateAsInt {
        result[index].DayTasks = append(groupedDayTask.DayTasks, dayTasks[i])
        found = true
      }
	  }
    if !found {
      groupedDayTask := domain.GroupedDayTask{ DateAsInt: dayTasks[i].DateAsInt, Display: common.ConvertIntToDisplay(dayTasks[i].DateAsInt), DayTasks: []domain.DayTask{ dayTasks[i] } }
      result = append(result, groupedDayTask)
    }
  }

  return result, err
}

func (r *TaskUseCases) FindTrashGroupedByDay(bookID string) ([]domain.GroupedDayTask, error) {
	result := []domain.GroupedDayTask{}

	// get a list of all the open tasks
	tasks, err := r.TaskRepository.Deleted(bookID)
	if err != nil {
    return result, err
  }

	// get the last open day from the last day ID for each task
	dayIds := []string{}
  for i := 0; i < len(tasks); i++ {
		log.Println("lastDayId: " + tasks[i].LastDayID)
    dayIds = append(dayIds, tasks[i].LastDayID)
  }
	days, err := r.DayRepository.FindByIds(bookID, dayIds)
	if err != nil {
    return result, err
  }

	// build a list of day tasks
	dayTasks := []domain.DayTask{}
	for i := 0; i < len(tasks); i++ {
    dayTask := domain.DayTask{}
    dayTask.Build(days[i], tasks[i])
		dayTasks = append(dayTasks, dayTask)
  }

	// sort the list by date descending
	sort.Sort(domain.DayTaskSort(dayTasks))

	// create a list of laterDays and add group the day tasks by unquie days
  for i := 0; i < len(dayTasks); i++ {
    // find the dateAsInt in the result
    found := false
    for index, groupedDayTask := range result {
      if groupedDayTask.DateAsInt == dayTasks[i].DateAsInt {
        result[index].DayTasks = append(groupedDayTask.DayTasks, dayTasks[i])
        found = true
      }
	  }
    if !found {
      groupedDayTask := domain.GroupedDayTask{ DateAsInt: dayTasks[i].DateAsInt, Display: common.ConvertIntToDisplay(dayTasks[i].DateAsInt), DayTasks: []domain.DayTask{ dayTasks[i] } }
      result = append(result, groupedDayTask)
    }
  }

  return result, err
}

func (r *TaskUseCases) NewTaskOnDate(bookID string, dayAsString string, elementName string) (domain.DayTask, error) {
	result := domain.DayTask{}

	dayAsInt, err := strconv.Atoi(dayAsString)
  if err != nil {
    return result, errors.New("Date is invalid")
  }

	// create the task, and add it to the requested date
  task, err := r.TaskRepository.Create(elementName, bookID, dayAsInt, "", r.User.ID)
  if err != nil {
    return result, err
  } 

  // create the linked day record
  day := domain.Day{BookID: bookID, TaskID: task.ID, Data: task.Data, DateAsInt: dayAsInt, IsActioned: false, IsCompleted: false, Created: time.Now()}
  day, err = r.DayRepository.Create(day)
	if err != nil {
    return result, err
  } 

  // build the result record
  result.Build(day, task)

  // return the new day task record
	return result, nil
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
	task.LastDayID = day.ID
	task.CurrentDate = 0
	task.Updated = time.Now()
	task.UpdatedBy = r.User.ID
	
	task, err = r.TaskRepository.Update(task)
	if err != nil {
		return result, err
	} 
  day, err = r.DayRepository.Update(day)
	result.Build(day, task)
  return result, err
}

func (r *TaskUseCases) DoOnDate(bookID string, taskId string, dateAsString string, data string) (domain.Task, domain.Day, error) {
	day := domain.Day{}

  // validate element and book
  task, err := r.TaskRepository.FindById(bookID, taskId)
  if err != nil || task.BookID != bookID {
		log.Println("BookID is not valid")
    return task, day, errors.New("BookID is not valid")
  }
  
  //if !task.LaterCanUpdate() {
	//	log.Println("Task cannot be updated")
  //  return task, day, errors.New("Task cannot be updated")
  //}

  convertedDate, err := common.ConvertStringToDates(dateAsString)
  if err != nil {
		log.Println("Date " + dateAsString + " is not valid")
    return task, day, errors.New("Date " + dateAsString + " is not valid")
  }
	log.Println(convertedDate)

  task.Data = data
  task.Updated = time.Now()
  task.UpdatedBy = r.User.ID
	task.IsDeleted = false

	if (task.CurrentDate != 0) {
		// close the current date
		log.Println("close the current day")
		currentDay, err := r.DayRepository.FindByTaskId(bookID, task.ID,  task.CurrentDate)
		currentDay.Data = data
		currentDay.IsActioned = true
		currentDay, err = r.DayRepository.Update(currentDay)
		if err != nil {
      return task, currentDay, err
    }
	}

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

func (r *TaskUseCases) EmptyTrash(bookID string) (error) {
	// delete all the trash items

	// get all the tasks 
	tasks, err := r.TaskRepository.Deleted(bookID)
	if err != nil {
		return nil
	}

	// for each task
	for i := 0; i < len(tasks); i++ {
		days, err := r.DayRepository.FindAllByTaskId(bookID, tasks[i].ID)
		if err == nil {
			for j := 0; j < len(days); j++ {
				// delete the day
				r.DayRepository.Delete(days[j].ID)
			}
			// delete the task
			r.TaskRepository.Delete(tasks[i].ID)
			log.Println("task deleted " + tasks[i].ID)

		} 
	}

	return nil

}




