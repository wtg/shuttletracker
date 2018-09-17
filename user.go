package shuttletracker

// User represents a user.
type User struct {
	ID       int64
	Username string
}

// UserService is an interface for interacting with Users.
type UserService interface {
	UserExists(username string) (bool, error)
	Users() ([]*User, error)
}
