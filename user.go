package shuttletracker

// User represents a user.
type User struct {
	Username string
}

type UserService interface {
	UserExists(username string) (bool, error)
}
