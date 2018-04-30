package handlers

import (
	"net/http"
	"strconv"
	"stufftobedone/templates"
)

//type activeTasks struct {
//}

/*ActiveTasks ... */
func ActiveTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dateAsString := query.Get("date")
	currentDateAsInt, err := strconv.Atoi(dateAsString)
	if err != nil {
		templates.JSONErrorMessage(w, "Date is invalid")
		return
	}

	c := getAppContext(r)
	// get the current date as int from the query parameters, otherwise default to server time
	defaultBook, err := c.BookUseCases.GetDefault(c.User)
	data, err := c.TaskUseCases.ActiveGroupedByDay(defaultBook.ID, currentDateAsInt)
	if err != nil {
		templates.JSONError(w, err)
		return
	}
	templates.JSON(w, data)
}
