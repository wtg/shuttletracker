package shuttletracker

import (
	"errors"
	"time"
)

// details of each form
type Form struct {
	ID      int64     `json:"id"`
	Message string    `json:"message"`
	Created time.Time `json:"created"`
	Admin	bool	  `json:"admin"`
}

// FeedbackService is an interface for interacting with Feedback.
type FeedbackService interface {
	Form(id int64) (*Form, error)
	Forms() ([]*Form, error)
	CreateForm(form *Form) error //idk if needs to be added with user input forms
	DeleteForm(id int64) error
}

// ErrFormNotFound indicates that a Form is not found.
var ErrFormNotFound = errors.New("Form not found")
