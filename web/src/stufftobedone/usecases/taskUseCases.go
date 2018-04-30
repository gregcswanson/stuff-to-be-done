package usecases

import (
	//"net/http"

	"net/http"
	"stufftobedone/domain"

	"github.com/gin-gonic/gin"
	//"appengine"
	//"appengine/datastore"
	"errors"
	"log"
	"sort"
	"strconv"
	"stufftobedone/common"
	"stufftobedone/repositories"
	"time"
)

type TaskUseCases struct {
	Request        *http.Request
	DayRepository  *repositories.DayRepository
	TaskRepository *repositories.TaskRepository
	BookRepository *repositories.BookRepository
	User           domain.AppUser
}

func NewTaskUseCases(c *gin.Context) *TaskUseCases {
	r := new(TaskUseCases)

	/*r.Request = request
	r.User = appUser
	r.DayRepository = repositories.NewDayRepository(r.Request)
	r.TaskRepository = repositories.NewTaskRepository(r.Request)*/
	return r
}

func NewTaskUseCases2(request *http.Request, appUser domain.AppUser) *TaskUseCases {
	r := new(TaskUseCases)
	r.Request = request
	r.User = appUser
	r.DayRepository = repositories.NewDayRepository(r.Request)
	r.TaskRepository = repositories.NewTaskRepository(r.Request)
	r.BookRepository = repositories.NewBookRepository(r.Request)
	return r
}

func (r *TaskUseCases) BuildDayTasks(days []domain.Day, tasks []domain.Task) []domain.DayTask {
	result := []domain.DayTask{}

	for i := 0; i < len(tasks); i++ {
		dayTask := domain.DayTask{}
		dayTask.Build(days[i], tasks[i])
		result = append(result, dayTask)
	}

	return result
}

func (r *TaskUseCases) FilterOutDeletedDayTasks(dayTasks []domain.DayTask) []domain.DayTask {
	result := []domain.DayTask{}

	for i := 0; i < len(dayTasks); i++ {
		if !dayTasks[i].IsDeleted {
			result = append(result, dayTasks[i])
		}
	}

	return result
}

func (r *TaskUseCases) ActiveGroupedByDay(bookID string, currentDateAsInt int) ([]domain.GroupedDayTask, error) {
	result := []domain.GroupedDayTask{}

	// get all item on the current day
	days, taskIds, err := r.DayRepository.DayWithTaskIds(bookID, currentDateAsInt)
	if err != nil {
		return result, err
	}

	tasks, err := r.TaskRepository.FindByIds(bookID, taskIds)
	if err != nil {
		return result, err
	}

	dayTasks := r.BuildDayTasks(days, tasks)
	dayTasks = r.FilterOutDeletedDayTasks(dayTasks)

	// create the current date group and fill it with tasks that have not been deleted
	currentDay := domain.GroupedDayTask{DateAsInt: currentDateAsInt, Display: common.ConvertIntToDisplay(currentDateAsInt), DayTasks: dayTasks}
	result = append(result, currentDay)

	// get the history
	days, taskIds, err = r.DayRepository.OpenItemsBeforeDayWithTaskIds(bookID, currentDateAsInt)
	if err != nil {
		return result, err
	}
	// get the tasks
	tasks, err = r.TaskRepository.FindByIds(bookID, taskIds)
	if err != nil {
		return result, err
	}
	// create a list of open daytasks
	openDayTasks := r.BuildDayTasks(days, tasks)
	// group them into the result
	for i := 0; i < len(openDayTasks); i++ {
		// find the dateAsInt in the result
		found := false
		for index, existingGroup := range result {
			if existingGroup.DateAsInt == openDayTasks[i].DateAsInt {
				result[index].DayTasks = append(existingGroup.DayTasks, openDayTasks[i])
				found = true
			}
		}
		if !found {
			newGroup := domain.GroupedDayTask{DateAsInt: openDayTasks[i].DateAsInt, Display: common.ConvertIntToDisplay(openDayTasks[i].DateAsInt), DayTasks: []domain.DayTask{openDayTasks[i]}}
			result = append(result, newGroup)
		}
	}

	return result, err
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
			groupedDayTask := domain.GroupedDayTask{DateAsInt: dayTasks[i].DateAsInt, Display: common.ConvertIntToDisplay(dayTasks[i].DateAsInt), DayTasks: []domain.DayTask{dayTasks[i]}}
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
			groupedDayTask := domain.GroupedDayTask{DateAsInt: dayTasks[i].DateAsInt, Display: common.ConvertIntToDisplay(dayTasks[i].DateAsInt), DayTasks: []domain.DayTask{dayTasks[i]}}
			result = append(result, groupedDayTask)
		}
	}

	return result, err
}

/*NewTask ... */
func (r *TaskUseCases) NewTask(dayTask domain.DayTask) (domain.DayTask, error) {
	result := domain.DayTask{}

	// create the task, and add it to the requested date
	task, err := r.TaskRepository.Create(dayTask.ElementName, dayTask.BookID, dayTask.DateAsInt, "", r.User.ID, dayTask.Data, dayTask.IsCompleted)
	if err != nil {
		return result, err
	}

	// create the linked to a day record
	day := domain.Day{BookID: dayTask.BookID, TaskID: task.ID, Data: task.Data, DateAsInt: dayTask.DateAsInt, IsActioned: false, IsCompleted: dayTask.IsCompleted, Created: time.Now()}
	day, err = r.DayRepository.Create(day)
	if err != nil {
		return result, err
	}

	// build the result record
	result.Build(day, task)

	// return the new day task record
	return result, nil
}

func (r *TaskUseCases) NewTaskOnDate(bookID string, dayAsString string, elementName string) (domain.DayTask, error) {
	result := domain.DayTask{}

	dayAsInt, err := strconv.Atoi(dayAsString)
	if err != nil {
		return result, errors.New("Date is invalid")
	}

	// create the task, and add it to the requested date
	task, err := r.TaskRepository.Create(elementName, bookID, dayAsInt, "", r.User.ID, "", false)
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

/*Update ... */
func (r *TaskUseCases) Update(dayTask domain.DayTask) (domain.DayTask, error) {
	// validate the book access?
	canAccess, err := r.BookRepository.ValidateAccess(dayTask.BookID, r.User.ID)
	if err != nil {
		return dayTask, err
	}
	if !canAccess {
		return dayTask, nil
	}

	// validate day and book
	task, err := r.TaskRepository.FindById(dayTask.BookID, dayTask.TaskID)
	if err != nil {
		return dayTask, err
	}

	// get the day item
	day, err := r.DayRepository.FindById(dayTask.BookID, dayTask.DayID)
	if err != nil {
		return dayTask, err
	}

	// update the details
	day.Data = dayTask.Data
	task.Data = dayTask.Data
	day.Tags = dayTask.Tags
	task.Tags = day.Tags
	task.Updated = time.Now()
	task.UpdatedBy = r.User.ID

	// if this day is the current day for the task, toggle the task coomplete
	if task.CurrentDate == day.DateAsInt {
		if !task.IsCompleted && dayTask.IsCompleted {
			task.Completed = time.Now()
		} else if task.IsCompleted && !dayTask.IsCompleted {
			task.Completed = time.Now()
		}
		task.IsCompleted = dayTask.IsCompleted
		day.IsCompleted = task.IsCompleted
	}

	// upate the repositories
	task, err = r.TaskRepository.Update(task, dayTask.BookID)
	if err != nil {
		return dayTask, err
	}
	day, err = r.DayRepository.Update(day)
	if err != nil {
		return dayTask, err
	}

	// build the new result record
	result := domain.DayTask{}
	result.Build(day, task)
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

		task, err = r.TaskRepository.Update(task, bookID)
		if err != nil {
			return result, err
		}
		result.BuildTask(task)
		return result, nil
	}

	// get the day item
	var day domain.Day
	if dayID == "" {
		day, err = r.DayRepository.FindByTaskId(bookID, task.ID, task.CurrentDate)
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

	task, err = r.TaskRepository.Update(task, bookID)
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

	if task.CurrentDate != 0 {
		// close the current date
		log.Println("close the current day")
		currentDay, err := r.DayRepository.FindByTaskId(bookID, task.ID, task.CurrentDate)
		currentDay.Data = data
		currentDay.IsActioned = true
		currentDay, err = r.DayRepository.Update(currentDay)
		if err != nil {
			return task, currentDay, err
		}
	}

	// create or find the day record
	day, err = r.DayRepository.FindByTaskId(bookID, task.ID, convertedDate.DateAsInt)
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
	task, err = r.TaskRepository.Update(task, bookID)
	return task, day, err
}

func (r *TaskUseCases) EmptyTrash(bookID string) error {
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
			r.TaskRepository.Delete(tasks[i].ID, bookID)
			log.Println("task deleted " + tasks[i].ID)

		}
	}

	return nil

}

func (r *TaskUseCases) SearchGroupedByDay(bookID string, search string) ([]domain.GroupedDayTask, error) {
	result := []domain.GroupedDayTask{}

	// search term cannot be blank
	if search == "" {
		return result, nil
	}

	// get a list of all the open tasks with the search term
	tasks, err := r.TaskRepository.Search(bookID, search)
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
			groupedDayTask := domain.GroupedDayTask{DateAsInt: dayTasks[i].DateAsInt, Display: common.ConvertIntToDisplay(dayTasks[i].DateAsInt), DayTasks: []domain.DayTask{dayTasks[i]}}
			result = append(result, groupedDayTask)
		}
	}

	return result, err
}
