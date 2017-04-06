package domain

import (
	"time"
)

// This struct is not stored in the datastore, it is used by the API only

type DayTask struct {
	DayID string
	TaskID string
	BookID string
	Data string
	ElementName string
	ProjectID string
	DateAsInt int
	IsCompleted bool
	IsActioned bool
	Created time.Time
	CanReopen bool
	Tags string
}

func (t *DayTask) Build(day Day, task Task)  {
	// merge the two entitied together
	t.DayID = day.ID
	t.TaskID = task.ID
	t.BookID = day.BookID
	t.Data = day.Data
	t.ElementName = task.ElementName
	t.ProjectID = task.ProjectID
	t.DateAsInt = day.DateAsInt
	t.IsCompleted = day.IsCompleted
	t.IsActioned = day.IsActioned
	t.Created = day.Created
	t.Tags = day.Tags
	if (task.CurrentDate == day.DateAsInt && t.IsCompleted){
		t.CanReopen = true
  	} else {
		t.CanReopen = false
	}
}

func (t *DayTask) BuildTask(task Task)  {
	// merge the two entitied together
	t.TaskID = task.ID
	t.Data = task.Data
	t.ElementName = task.ElementName
	t.ProjectID = task.ProjectID
	t.DateAsInt = 0
	t.IsCompleted = task.IsCompleted
	t.Created = task.Created
	t.Tags = task.Tags
	if (task.CurrentDate == 0 && t.IsCompleted){
		t.CanReopen = true
  	} else {
		t.CanReopen = false
	}
}

