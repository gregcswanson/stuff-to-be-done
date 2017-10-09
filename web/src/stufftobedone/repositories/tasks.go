package repositories

import (
	"net/http"
    "stufftobedone/domain"
	"appengine"
	"appengine/datastore"
    "time"
    "log"
    "strconv"
)

type TaskRepository struct {
  request    *http.Request
  cache      *AppCache
}

func NewTaskRepository(request *http.Request) *TaskRepository {
	r := new(TaskRepository)
	r.request = request
    r.cache = NewAppCache(request)
	return r
}

func (r *TaskRepository) Later(bookID string) ([]domain.Task, error) {
    tasks := []domain.Task{}
    globalContext := appengine.NewContext(r.request)
    
    // filter Tasks where not current scheduled, is not part of a project and is not completed
    q := datastore.NewQuery("Tasks").Filter("CurrentDate = ", 0).Filter("ProjectID =", "").Filter("IsCompleted =", false).Filter("IsDeleted =", false).Order("Created")

    // link the task to the book
    bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
    q = q.Ancestor(bookKey)

    keys, err := q.GetAll(globalContext, &tasks)
    if err != nil {    
        return tasks, err
    } else {    
        // loop through and add the keys as ID
        for i := 0; i < len(keys); i++ {
            tasks[i].ID = keys[i].Encode()
        }
    }
    
    return tasks, nil
}

func (r *TaskRepository) LaterCount(bookID string)(int, error) {
    // check for cached result
    cacheValue, foundInCache := r.cache.GetInt("latercount", bookID)
    if foundInCache {
        return cacheValue, nil
    }
    laterItems, err := r.Later(bookID)
    if err != nil {    
        return 0, err
    }
    r.cache.SetInt("latercount", bookID, len(laterItems))
    
    return len(laterItems), nil
}

func (r *TaskRepository) Completed(bookID string) ([]domain.Task, error) {
    tasks := []domain.Task{}
    globalContext := appengine.NewContext(r.request)
    
    // filter Tasks where not current scheduled, is not part of a project and is completed
    q := datastore.NewQuery("Tasks").Filter("CurrentDate = ", 0).Filter("ProjectID =", "").Filter("IsCompleted =", true).Filter("IsDeleted =", false).Order("-Completed")

    // link the task to the book
    bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
    q = q.Ancestor(bookKey)

    keys, err := q.GetAll(globalContext, &tasks)
    if err != nil {  
        log.Println("No Tasks found", err)  
        return tasks, err
    } else {    
        // loop through and add the keys as ID
        for i := 0; i < len(keys); i++ {
            tasks[i].ID = keys[i].Encode()
        }
    }
    
    return tasks, nil
}

func (r *TaskRepository) Deleted(bookID string) ([]domain.Task, error) {
    tasks := []domain.Task{}
    globalContext := appengine.NewContext(r.request)
    
    // filter Tasks where not current scheduled, is not part of a project and is completed
    q := datastore.NewQuery("Tasks").Filter("CurrentDate = ", 0).Filter("ProjectID =", "").Filter("IsDeleted =", true).Order("-Updated")

    // link the task to the book
    bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
    q = q.Ancestor(bookKey)

    keys, err := q.GetAll(globalContext, &tasks)
    if err != nil {  
        log.Println("No Tasks found", err)  
        return tasks, err
    } else {    
        log.Println("Tasks found")
        // loop through and add the keys as ID
        for i := 0; i < len(keys); i++ {
            tasks[i].ID = keys[i].Encode()
        }
    }
    
    return tasks, nil
}

func (r *TaskRepository) FindByIds(bookId string, ids []string) ([]domain.Task, error){
     // create the namespace context
    globalContext := appengine.NewContext(r.request)
   
    // convert the ids to keys
    length := len(ids)
    tasks := make([]domain.Task, length)
    keys := []*datastore.Key{}
    for i := 0; i < length; i++ {
        key, _ := datastore.DecodeKey(ids[i])
        keys = append(keys, key)
    }
    err := datastore.GetMulti(globalContext, keys, tasks)
    if err == nil {
        for i := 0; i < length; i++ {
            tasks[i].ID = ids[i]
        }
    } 
    return tasks, err

}

func (r *TaskRepository) FindById(bookID string, taskID string) (domain.Task, error){
    var task domain.Task

    // create the namespace context
    globalContext := appengine.NewContext(r.request)

    // link the task to the book
    // bookKey := datastore.NewKey(globalContext, "Tasks", bookID, 0, nil)
    // q = q.Ancestor(bookKey)

    // get the key
    key, err := datastore.DecodeKey(taskID)
    if err != nil {
        return task, err
    }
    // retrieve the task
    err = datastore.Get(globalContext, key, &task);
    task.ID = taskID

    return task, err
}

func (r *TaskRepository) Create(elementName string, bookID string, dayAsInt int, projectID string, userID string) (domain.Task, error) {
    // clear the count cache
    r.cache.Clear("latercount", bookID)

    // to-do validate book, element name and projectID

    // to-do if the day is not null add to the day reference table also
    // to-do convert to UTC now
    task := domain.Task{ ID: "", BookID: bookID, Data: "", ElementName: elementName, ProjectID: projectID, IsCompleted: false, CurrentDate: dayAsInt, Created: time.Now(), CreatedBy: userID, Updated: time.Now(), UpdatedBy: userID }
    task.DueDateAsInt = 99999999
    task.Sort, _ = strconv.Atoi(task.Created.Format("20060102150405"))
    globalContext := appengine.NewContext(r.request)
    bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
    
  	err := datastore.RunInTransaction(globalContext, func(c appengine.Context) error {
        
        key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Tasks", bookKey), &task)
        if err != nil {
            return err
        } else {
            task.ID = key.Encode()
        }
        
        return nil
    }, nil)
    if err != nil {
        return task, err
    }

    return task, nil
}

func (r *TaskRepository) Update(task domain.Task, bookID string) (domain.Task, error){
    // clear the count cache
    r.cache.Clear("latercount", bookID)

    globalContext := appengine.NewContext(r.request)
    key, err := datastore.DecodeKey(task.ID)
    if err != nil {
        return task, err
    }
    _, err = datastore.Put(globalContext, key, &task)
    return task, err
}

func (r *TaskRepository) Delete(taskID string, bookID string) (error){
    // clear the count cache
    r.cache.Clear("latercount", bookID)

    globalContext := appengine.NewContext(r.request)
  
    key , err := datastore.DecodeKey(taskID)
	if err != nil {
		return err
	}
	err = datastore.Delete(globalContext, key)
    return err
}
// - set default value, save and return with ID
// update
// - update task and day table
// - task archive table?

// delete

// delete from day only - needs a option on the container element for this