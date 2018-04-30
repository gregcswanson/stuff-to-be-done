package handlers

import (
	"encoding/json"
	"net/http"
	"stufftobedone/domain"
	"stufftobedone/templates"
)

/*PostDayTasks ... */
func PostDayTasks(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)

	// read the task details from the body
	var dayTask domain.DayTask
	err := json.NewDecoder(r.Body).Decode(&dayTask)
	if err != nil {
		templates.JSONError(w, err)
		return
	}

	createdDayTask, err := c.TaskUseCases.NewTask(dayTask)
	if err != nil {
		templates.JSONError(w, err)
		return
	}
	createdDayTask.ClientID = dayTask.ClientID
	templates.JSON(w, createdDayTask)
}
