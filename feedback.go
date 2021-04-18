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
	Admin	bool	  `json:"admin_shuttletracker"`
}

// FeedbackService is an interface for interacting with Feedback.
type FeedbackService interface {
	GetAdminForm() (*Form)
	GetForm(id int64) (*Form, error)
	GetForms() ([]*Form, error)
	CreateForm(form *Form) error
	DeleteForm(id int64) error
}

// ErrFormNotFound indicates that a Form is not found.
var ErrFormNotFound = errors.New("Form not found")
