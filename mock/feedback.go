package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

// FeedbackService implements a mock of shuttletracker.FeedbackService.
type FeedbackService struct {
	mock.Mock
}

// CreateForm creates a Form.
func (fs *FeedbackService) CreateForm(form *shuttletracker.Form) error {
	args := fs.Called(form)
	return args.Error(0)
}

// DeleteForm deletes a form
func (fs *FeedbackService) DeleteForm(id int64) error {
	args := fs.Called(id)
	return args.Error(0)
}

// Form gets a form
func (fs *FeedbackService) Form(id int64) (*shuttletracker.Form, error) {
	args := fs.Called(id)
	return args.Get(0).(*shuttletracker.Form), args.Error(1)
}

// EditForm edits form
func (fs *FeedbackService) EditForm(form *shuttletracker.Form) error {
	args := fs.Called(route)
	return args.Error(0)
}

// Forms returns all forms
func (fs *FeedbackService) Forms() ([]*shuttletracker.Form, error) {
	args := fs.Called()
	return args.Get(0).([]*shuttletracker.Form), args.Error(1)
}
