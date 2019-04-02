package shuttletracker

import "errors"

// ErrUserNotFound indicates that a User is not in the service.
var ErrUserNotFound = errors.New("User not found")

// User represents a user.
type User struct {
	ID       int64
	Username string
}

// UserService is an interface for interacting with Users.
type UserService interface {
	CreateUser(*User) error
	DeleteUser(username string) error
	UserExists(username string) (bool, error)
	Users() ([]*User, error)
}
