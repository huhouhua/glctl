package progress

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"strings"
	"time"
)

type EventStatus int

const (
	Running EventStatus = iota
	// Done means that the current task is done
	Done
	// Warning means that the current task has warning
	Warning
	// Error means that the current task has errored
	Error
)

// Event represents a progress event.
type Event struct {
	isParent bool
	status   EventStatus
	color    string
	spinner  *spinner.Spinner
}

// ErrorEvent creates a new Error Event
func CreatingEvent(isParent bool) *Event {
	return NewEvent(Running, isParent, "green")
}

// NewEvent new event
func NewEvent(status EventStatus, isParent bool, colors ...string) *Event {
	return &Event{
		status:   status,
		isParent: isParent,
		spinner: spinner.New(spinner.CharSets[11], 100*time.Millisecond, func(s *spinner.Spinner) {
			if !isParent {
				s.Prefix = " "
			}
			_ = s.Color(colors...)
		}),
	}
}
func (e *Event) WithText(text string) *Event {
	e.spinner.Suffix = text
	return e
}

func (e *Event) Start() *Event {
	e.spinner.Start()
	return e
}
func (e *Event) Stop() *Event {
	e.spinner.Stop()
	return e
}
func (e *Event) Success(msg ...string) *Event {
	e.status = Done
	e.setFinal(msg...)
	return e
}

func (e *Event) Done(msg ...string) *Event {
	e.status = Done
	e.spinner.FinalMSG = fmt.Sprintf("%s %s\n", strings.Join(msg, ""), color.GreenString("successfully"))
	return e
}

func (e *Event) Error(msg ...string) *Event {
	e.status = Error
	e.setFinal(msg...)
	return e
}

func (e *Event) setFinal(msg ...string) {
	s := e.Spinner()
	if !e.isParent {
		s = " " + s
	}
	e.spinner.FinalMSG = fmt.Sprintf("%s %s %s\n", s, e.spinner.Suffix, strings.Join(msg, " "))
}

var (
	spinnerDone    = "✔"
	spinnerWarning = "!"
	spinnerError   = "✘"
)

func (e *Event) Spinner() string {
	switch e.status {
	case Done:
		return SuccessColor(spinnerDone)
	case Warning:
		return WarningColor(spinnerWarning)
	case Error:
		return ErrorColor(spinnerError)
	case Running:
		return ""
	default:
		return CountColor(e.spinner.Suffix)
	}
}
