package progress

import (
	"github.com/briandowns/spinner"
	"time"
)

type EventStatus int

const (
	// Done means that the current task is done
	Done EventStatus = iota
	// Warning means that the current task has warning
	Warning
	// Error means that the current task has errored
	Error
)

// Event represents a progress event.
type Event struct {
	ID         string
	ParentID   string
	Text       string
	Status     EventStatus
	StatusText string
	Current    int64
	Percent    int

	Total     int64
	startTime time.Time
	endTime   time.Time
	spinner   *spinner.Spinner
}

// ErrorMessageEvent creates a new Error Event with message
func ErrorMessageEvent(id string, msg string) Event {
	return NewEvent(id, Error, msg)
}

// NewEvent new event
func NewEvent(id string, status EventStatus, statusText string) Event {
	return Event{
		ID:         id,
		Status:     status,
		StatusText: statusText,
	}
}
