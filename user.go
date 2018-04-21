package shuttletracker

// User represents a user.
type User struct {
	Username string
}

// UserService is an interface for interacting with Users.
type UserService interface {
	UserExists(username string) (bool, error)
}
