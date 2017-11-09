package domain

import (
	"time"
)

// This struct is not stored in the datastore, it is used by the API only

type DayTask struct {
	DayID       string
	TaskID      string
	BookID      string
	Data        string
	ElementName string
	ProjectID   string
	Search		string
	DateAsInt   int
	IsCompleted bool
	IsActioned  bool
	Created     time.Time
	CanReopen   bool
	Tags        string
	Comment     string // a comment for the day
}

func (t *DayTask) Build(day Day, task Task) {
	// merge the two entitied together
	t.DayID = day.ID
	t.TaskID = task.ID
	t.BookID = day.BookID
	t.Data = day.Data
	t.ElementName = task.ElementName
	t.ProjectID = task.ProjectID
	t.Search = task.Search
	t.DateAsInt = day.DateAsInt
	t.IsCompleted = day.IsCompleted
	t.IsActioned = day.IsActioned
	t.Created = day.Created
	t.Tags = day.Tags
	t.Comment = day.Comment
	if task.CurrentDate == day.DateAsInt && t.IsCompleted {
		t.CanReopen = true
	} else {
		t.CanReopen = false
	}
}

func (t *DayTask) BuildTask(task Task) {
	// merge the two entitied together
	t.TaskID = task.ID
	t.Data = task.Data
	t.ElementName = task.ElementName
	t.ProjectID = task.ProjectID
	t.Search = task.Search
	t.DateAsInt = 0
	t.IsCompleted = task.IsCompleted
	t.Created = task.Created
	t.Tags = task.Tags
	t.Comment = ""
	if task.CurrentDate == 0 && t.IsCompleted {
		t.CanReopen = true
	} else {
		t.CanReopen = false
	}
}

type DayTaskSort []DayTask

func (a DayTaskSort) Len() int           { return len(a) }
func (a DayTaskSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DayTaskSort) Less(i, j int) bool { return a[i].DateAsInt > a[j].DateAsInt }
