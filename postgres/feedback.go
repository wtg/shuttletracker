package postgres

import (
	"database/sql"

	"github.com/wtg/shuttletracker"
)

// FeedbackService is an implementation of shuttletracker.FeedbackService.
type FeedbackService struct {
	db *sql.DB
}

func (fs *FeedbackService) initializeSchema(db *sql.DB) error {
	fs.db = db
	schema := `
CREATE TABLE IF NOT EXISTS forms (
	id serial PRIMARY KEY,
	topic text,
	message text,
	created timestamp with time zone NOT NULL DEFAULT now(),
	read bool NOT NULL
);`
	_, err := fs.db.Exec(schema)
	return err
}

// Form returns a Form by its ID.
func (fs *FeedbackService) Form(id int64) (*shuttletracker.Form, error) {
	f := &shuttletracker.Form{
		ID: id,
	}

	statement := "SELECT f.created, f.topic, f.message, f.read" +
		" FROM forms f WHERE id = $1;"
	row := fs.db.QueryRow(statement, id)
	err := row.Scan(&f.Created, &f.Topic, &f.Message, &f.Read)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrFormNotFound
	}

	return f, err
}

// Forms returns all Forms.
func (fs *FeedbackService) Forms() ([]*shuttletracker.Form, error) {
	forms := []*shuttletracker.Form{}
	query := "SELECT f.id, f.topic, f.created, f.message, f.read" +
		" FROM forms f;"
	rows, err := fs.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		f := &shuttletracker.Form{}
		err := rows.Scan(&f.ID, &f.Topic, &f.Created, &f.Message, &f.Read)
		if err != nil {
			return nil, err
		}
		forms = append(forms, f)
	}
	return forms, nil
}

// not sure if properly made
// EditForm updates read status of the form
func (fs *FeedbackService) EditForm(form *shuttletracker.Form) error {
	tx, err := fs.db.Begin()
	if err != nil {
		return err
	}
	// We can't really do anything if rolling back a transaction fails.
	// nolint: errcheck
	defer tx.Rollback()

	// change read status
	statement := "UPDATE forms SET read = $1" +
		" WHERE id = $2 RETURNING read;"
	row := tx.QueryRow(statement, form.Read, form.ID)
	// not sure if able to scan like this;
	// wanted to return a true/false to show it was successfully set
	err = row.Scan(&form.Read)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// CreateForm creates a Form.
func (fs *FeedbackService) CreateForm(id int64) error {
	statement := "INSERT INTO forms (topic, message) VALUES" +
		" ($1, $2) RETURNING id, created, read;"
	f := &shuttletracker.Form{}
	row := fs.db.QueryRow(statement, f.Topic, f.Message)
	return row.Scan(&f.ID, &f.Created, &f.Read)
}

// DeleteForm deletes a Form.
func (fs *FeedbackService) DeleteForm(id int64) error {
	statement := "DELETE FROM forms WHERE id = $1;"
	result, err := fs.db.Exec(statement, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return shuttletracker.ErrFormNotFound
	}

	return nil
}
