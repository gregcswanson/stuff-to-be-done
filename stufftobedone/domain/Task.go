package domain

import (
	"time"
)

type Task struct {
	ID string `datastore:"-"`
	BookID string // link to the payments preocessor for multi user books
	Data string `datastore:",noindex"`
	ElementName string
	ProjectID string
	CurrentDate int // the most recent date this task is assigned to
	IsComplete bool
	IsDeleted bool
	DueDate time.Time `datastore:"-"`
	DueDateAsInt int // store the date as int for sorting, no date is 99999999 so that it is last
	Created time.Time
	CreatedBy string `datastore:",noindex"`
  	Updated time.Time `datastore:",noindex"` // future enhancement to track time spent
	UpdatedBy string `datastore:",noindex"`
  	Completed time.Time //`datastore:",noindex"` // future enhancement to track time spent
	CompletedBy string `datastore:",noindex"`
	Tags string
}

func (t Task) LaterCanUpdate() bool {
	// validate this task can be edited by the later command
	if t.IsComplete && t.ProjectID == "" {
		return true
	}
	if !t.IsComplete && t.CurrentDate == 0 && t.ProjectID == "" {
		return true
	}
	return false
}

func (t Task) LaterCanDelete() bool {
	// validate this task can be edited by the later command
	if t.IsComplete && t.ProjectID == "" {
		return true
	}
	if !t.IsComplete && t.CurrentDate == 0 && t.ProjectID == "" {
		return true
	}
	return false
}

func (t Task) LaterCanUndoDelete() bool {
	// validate this task can be edited by the later command
	if t.IsComplete && t.ProjectID == "" && t.IsDeleted {
		return true
	}
	if !t.IsComplete && t.CurrentDate == 0 && t.ProjectID == ""&& t.IsDeleted {
		return true
	}
	return false
}

