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
	message text,
	created timestamp with time zone NOT NULL DEFAULT now(),
	read bool NOT NULL,
);`
	_, err := fs.db.Exec(schema)
	return err
}

// not sure if needs to be added ! check
// CreateForm creates a Form.
func (fs *FeedbackService) CreateForm(form *shuttletracker.Form) error {
	statement := "INSERT INTO forms (message) VALUES" +
		" ($1) RETURNING id, created, read;"
	row := fs.db.QueryRow(statement, form.Message)
	return row.Scan(&form.ID, &form.Created, &form.Read)
}

// Form returns a Form by its ID.
func (fs *FeedbackService) Form(id int64) (*shuttletracker.Form, error) {
	f := &shuttletracker.Form{
		ID: id,
	}

	statement := "SELECT f.created, f.message, f.read" +
		" FROM forms f WHERE id = $1;"
	row := fs.db.QueryRow(statement, id)
	err := row.Scan(&f.Created, &f.Message, &f.Read)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrFormNotFound
	}

	return f, err
}

// Forms returns all Forms.
func (fs *FeedbackService) Forms() ([]*shuttletracker.Form, error) {
	forms := []*shuttletracker.Form{}
	query := "SELECT f.id, f.created, f.message, f.read" +
		" FROM forms s;"
	rows, err := fs.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		f := &shuttletracker.Form{}
		err := rows.Scan(&f.ID, &f.Created, &f.Message, &f.Read)
		if err != nil {
			return nil, err
		}
		forms = append(forms, f)
	}
	return forms, nil
}

// DeleteStop deletes a Stop.
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
