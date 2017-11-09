package domain

import (
	"time"
)

type Task struct {
	ID           string `datastore:"-"`
	BookID       string // link to the payments preocessor for multi user books
	Data         string `datastore:",noindex"`
	ElementName  string
	ProjectID    string
	Search		 string
	Queue        string
	Priority     int
	ParentTaskID string // The task this record was triggered from, i.e. part of meeting turned into a todo
	CurrentDate  int    // the most recent date this task is assigned to
	LastDayID    string // ths last day this task was attached to
	IsCompleted  bool
	IsDeleted    bool
	DueDate      time.Time `datastore:"-"`
	DueDateAsInt int       // store the date as int for sorting, no date is 99999999 so that it is last
	Sort         int
	Created      time.Time
	CreatedBy    string    `datastore:",noindex"`
	Updated      time.Time //`datastore:",noindex"`
	UpdatedBy    string    `datastore:",noindex"`
	Completed    time.Time //`datastore:",noindex"`
	CompletedBy  string    `datastore:",noindex"`
	Tags         string
}

func (t Task) LaterCanUpdate() bool {
	// validate this task can be edited by the later command
	if t.IsCompleted && t.ProjectID == "" {
		return true
	}
	if !t.IsCompleted && t.CurrentDate == 0 && t.ProjectID == "" {
		return true
	}
	return false
}

func (t Task) LaterCanDelete() bool {
	// validate this task can be edited by the later command
	if t.IsCompleted && t.ProjectID == "" {
		return true
	}
	if !t.IsCompleted && t.CurrentDate == 0 && t.ProjectID == "" {
		return true
	}
	return false
}

func (t Task) LaterCanUndoDelete() bool {
	// validate this task can be edited by the later command
	if t.IsCompleted && t.ProjectID == "" && t.IsDeleted {
		return true
	}
	if !t.IsCompleted && t.CurrentDate == 0 && t.ProjectID == "" && t.IsDeleted {
		return true
	}
	return false
}
