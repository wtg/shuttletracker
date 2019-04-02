package postgres

import (
	"testing"

	"github.com/wtg/shuttletracker"
)

// nolint: gocyclo
func TestCreateDeleteUser(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	const username = "testuser"

	exists, err := pg.UserExists(username)
	if err != nil {
		t.Fatalf("unable to check if user exists: %s", err)
	}
	if exists {
		t.Fatalf("user already exists")
	}

	user := &shuttletracker.User{
		Username: username,
	}
	err = pg.CreateUser(user)
	if err != nil {
		t.Fatalf("unable to create User: %s", err)
	}

	exists, err = pg.UserExists(username)
	if err != nil {
		t.Fatalf("unable to check if user exists: %s", err)
	}
	if !exists {
		t.Fatalf("user does not exist")
	}

	err = pg.DeleteUser(username)
	if err != nil {
		t.Fatalf("unable to delete User: %s", err)
	}

	exists, err = pg.UserExists(username)
	if err != nil {
		t.Fatalf("unable to check if user exists: %s", err)
	}
	if exists {
		t.Fatalf("user still exists")
	}
}

func TestErrUserNotFound(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	const username = "testuser"

	err := pg.DeleteUser(username)
	if err != shuttletracker.ErrUserNotFound {
		t.Errorf("got unexpected error: %s", err)
	}
}

// nolint: gocyclo
func TestUsers(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	const username1 = "testuser1"
	const username2 = "testuser2"

	user := &shuttletracker.User{
		Username: username1,
	}
	err := pg.CreateUser(user)
	if err != nil {
		t.Fatalf("unable to create User: %s", err)
	}

	user = &shuttletracker.User{
		Username: username2,
	}
	err = pg.CreateUser(user)
	if err != nil {
		t.Fatalf("unable to create User: %s", err)
	}

	users, err := pg.Users()
	if err != nil {
		t.Fatalf("unable to get users: %s", err)
	}

	if len(users) != 2 {
		t.Fatalf("got %d users, expected 2", len(users))
	}

	user1 := false
	user2 := false
	for _, user := range users {
		if user.Username == username1 {
			user1 = true
		} else if user.Username == username2 {
			user2 = true
		}
	}

	if !user1 || !user2 {
		t.Fatalf("not all users returned")
	}
}
