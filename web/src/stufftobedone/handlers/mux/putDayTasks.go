package handlers

import (
	"encoding/json"
	"net/http"
	"stufftobedone/domain"
	"stufftobedone/templates"
)

/*PutDayTasks ... */
func PutDayTasks(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)

	// read the task details from the body
	var dayTask domain.DayTask
	err := json.NewDecoder(r.Body).Decode(&dayTask)
	if err != nil {
		templates.JSONError(w, err)
		return
	}

	updatedDayTask, err := c.TaskUseCases.Update(dayTask)
	if err != nil {
		templates.JSONError(w, err)
		return
	}
	updatedDayTask.ClientID = dayTask.ClientID
	templates.JSON(w, updatedDayTask)
}
