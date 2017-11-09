package domain

import (
	"time"
)

type Day struct {
	ID          string `datastore:"-"`
	BookID      string // link to the payments preocessor for multi user books
	TaskID      string
	Data        string `datastore:",noindex"` // saved for versioning
	DateAsInt   int
	Sort        int
	Comment     string // a comment for the day
	IsCompleted bool
	IsActioned  bool // was this actioned, sent to later, completed, the display will be the read-only version
	Created     time.Time
	Tags        string
}
