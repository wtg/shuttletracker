package postgres

import (
	"database/sql"

	"github.com/wtg/shuttletracker"
)

// UserService is an implementation of shuttletracker.UserService.
type UserService struct {
	db *sql.DB
}

func (us *UserService) initializeSchema(db *sql.DB) error {
	us.db = db
	schema := `
CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	username varchar(10) UNIQUE NOT NULL
);
	`
	_, err := us.db.Exec(schema)
	return err
}

// UserExists returns whether a User with the specified username exists.
func (us *UserService) UserExists(username string) (bool, error) {
	row := us.db.QueryRow("SELECT FROM users WHERE username = $1;", username)
	err := row.Scan()
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Users returns all existing Users..
func (us *UserService) Users() ([]*shuttletracker.User, error) {
	var users []*shuttletracker.User

	statement := "SELECT id, username FROM users;"
	rows, err := us.db.Query(statement)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := &shuttletracker.User{}
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}
