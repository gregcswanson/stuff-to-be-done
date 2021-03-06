package repositories

import (
	"errors"
	"net/http"
	"strconv"
	"stufftobedone/domain"

	"appengine"
	"appengine/datastore"
)

type DayRepository struct {
	request *http.Request
}

func NewDayRepository(request *http.Request) *DayRepository {
	r := new(DayRepository)
	r.request = request
	return r
}

func (r *DayRepository) Day(bookID string, dateAsInt int) ([]domain.Day, error) {
	days := []domain.Day{}
	globalContext := appengine.NewContext(r.request)

	// filter Days where not current scheduled, is not part of a project and is not completed
	//q := datastore.NewQuery("Days").Filter("DateAsInt = ", dateAsInt).Filter("IsActioned = ", false).Order("Created")
	q := datastore.NewQuery("Days").Filter("DateAsInt = ", dateAsInt).Order("Created")

	// link the day to the book
	bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
	q = q.Ancestor(bookKey)

	keys, err := q.GetAll(globalContext, &days)
	if err != nil {
		return days, err
	} else {
		// loop through and add the keys as ID
		for i := 0; i < len(keys); i++ {
			days[i].ID = keys[i].Encode()
		}
	}

	return days, nil
}

func (r *DayRepository) DayWithTaskIds(bookID string, dateAsInt int) ([]domain.Day, []string, error) {
	days := []domain.Day{}
	taskIds := []string{}
	globalContext := appengine.NewContext(r.request)

	// filter Days where not current scheduled, is not part of a project and is not completed
	//q := datastore.NewQuery("Days").Filter("DateAsInt = ", dateAsInt).Filter("IsActioned = ", false).Order("Created")
	q := datastore.NewQuery("Days").Filter("DateAsInt = ", dateAsInt).Order("-Created")

	// link the day to the book
	bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
	q = q.Ancestor(bookKey)

	keys, err := q.GetAll(globalContext, &days)
	if err != nil {
		return days, taskIds, err
	} else {
		// loop through and add the keys as ID
		for i := 0; i < len(keys); i++ {
			days[i].ID = keys[i].Encode()
		}
	}

	// read the task ids into a list
	for i := 0; i < len(days); i++ {
		taskIds = append(taskIds, days[i].TaskID)
	}

	return days, taskIds, nil
}

func (r *DayRepository) OpenItemsBeforeDayWithTaskIds(bookID string, dateAsInt int) ([]domain.Day, []string, error) {
	days := []domain.Day{}
	taskIds := []string{}
	globalContext := appengine.NewContext(r.request)

	// filter Days where not current scheduled, is not part of a project and is not completed
	q := datastore.NewQuery("Days").Filter("DateAsInt < ", dateAsInt).Filter("IsActioned = ", false).Filter("IsCompleted = ", false).Order("-DateAsInt").Order("-Created")

	// link the day to the book
	bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
	q = q.Ancestor(bookKey)

	keys, err := q.GetAll(globalContext, &days)
	if err != nil {
		return days, taskIds, err
	} else {
		// loop through and add the keys as ID
		for i := 0; i < len(keys); i++ {
			days[i].ID = keys[i].Encode()
		}
	}

	// read the task ids into a list
	for i := 0; i < len(days); i++ {
		taskIds = append(taskIds, days[i].TaskID)
	}

	return days, taskIds, nil
}

func (r *DayRepository) OpenItemsBeforeDay(bookID string, dateAsInt int) ([]domain.Day, error) {
	days := []domain.Day{}
	globalContext := appengine.NewContext(r.request)

	// filter Days where not current scheduled, is not part of a project and is not completed
	q := datastore.NewQuery("Days").Filter("DateAsInt < ", dateAsInt).Filter("IsActioned = ", false).Filter("IsCompleted = ", false).Order("-DateAsInt").Order("Created")

	// link the day to the book
	bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
	q = q.Ancestor(bookKey)

	keys, err := q.GetAll(globalContext, &days)
	if err != nil {
		return days, err
	} else {
		// loop through and add the keys as ID
		for i := 0; i < len(keys); i++ {
			days[i].ID = keys[i].Encode()
		}
	}

	return days, nil
}
func (r *DayRepository) FindById(bookID string, dayID string) (domain.Day, error) {
	var day domain.Day

	// create the namespace context
	globalContext := appengine.NewContext(r.request)

	// get the key
	key, err := datastore.DecodeKey(dayID)
	if err != nil {
		return day, err
	}
	// retrieve the day
	err = datastore.Get(globalContext, key, &day)
	day.ID = dayID

	return day, err
}

func (r *DayRepository) FindByTaskId(bookID string, taskID string, dateAsInt int) (domain.Day, error) {
	var day domain.Day
	days := []domain.Day{}

	// create the namespace context
	globalContext := appengine.NewContext(r.request)

	q := datastore.NewQuery("Days").Filter("DateAsInt = ", dateAsInt).Filter("TaskID = ", taskID)

	// link the day to the book
	bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
	q = q.Ancestor(bookKey)

	keys, err := q.GetAll(globalContext, &days)
	if err != nil {
		return day, err
	}
	if len(keys) == 0 {
		return day, errors.New("Task not found")
	}
	day = days[0]
	day.ID = keys[0].Encode()
	return day, nil
}

func (r *DayRepository) FindAllByTaskId(bookID string, taskID string) ([]domain.Day, error) {
	days := []domain.Day{}
	globalContext := appengine.NewContext(r.request)

	// filter Days where not current scheduled, is not part of a project and is not completed
	q := datastore.NewQuery("Days").Filter("TaskID = ", taskID)

	// link the day to the book
	bookKey := datastore.NewKey(globalContext, "Books", bookID, 0, nil)
	q = q.Ancestor(bookKey)

	keys, err := q.GetAll(globalContext, &days)
	if err != nil {
		return days, err
	} else {
		// loop through and add the keys as ID
		for i := 0; i < len(keys); i++ {
			days[i].ID = keys[i].Encode()
		}
	}

	return days, nil
}

func (r *DayRepository) FindByIds(bookId string, ids []string) ([]domain.Day, error) {
	// create the namespace context
	globalContext := appengine.NewContext(r.request)

	// convert the ids to keys
	length := len(ids)
	days := make([]domain.Day, length)
	keys := []*datastore.Key{}
	for i := 0; i < length; i++ {
		key, _ := datastore.DecodeKey(ids[i])
		keys = append(keys, key)
	}
	err := datastore.GetMulti(globalContext, keys, days)
	if err == nil {
		for i := 0; i < length; i++ {
			days[i].ID = ids[i]
		}
	}
	return days, err

}

func (r *DayRepository) Create(day domain.Day) (domain.Day, error) {

	globalContext := appengine.NewContext(r.request)
	bookKey := datastore.NewKey(globalContext, "Books", day.BookID, 0, nil)
	day.Sort, _ = strconv.Atoi(day.Created.Format("20060102150405"))
	err := datastore.RunInTransaction(globalContext, func(c appengine.Context) error {

		key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Days", bookKey), &day)
		if err != nil {
			return err
		} else {
			day.ID = key.Encode()
		}

		return nil
	}, nil)
	if err != nil {
		return day, err
	}

	return day, nil
}

func (r *DayRepository) Update(day domain.Day) (domain.Day, error) {
	globalContext := appengine.NewContext(r.request)
	key, err := datastore.DecodeKey(day.ID)
	if err != nil {
		return day, err
	}
	_, err = datastore.Put(globalContext, key, &day)
	return day, err
}

func (r *DayRepository) Delete(dayID string) error {
	globalContext := appengine.NewContext(r.request)

	key, err := datastore.DecodeKey(dayID)
	if err != nil {
		return err
	}
	err = datastore.Delete(globalContext, key)
	return err
}
