package postgres

import "database/sql"

// UserService is an implementation of shuttletracker.UserService.
type UserService struct {
	db *sql.DB
}

func (us *UserService) initializeSchema() error {
	schema := ``
	_, err := us.db.Exec(schema)
	return err
}

// UserExists returns whether a User with the specified username exists.
func (us *UserService) UserExists(username string) (bool, error) {
	return false, nil
}
