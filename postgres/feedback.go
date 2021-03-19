package postgres

import (
	"database/sql"
	"strconv"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
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
	prompt text,
	message text,
	created timestamp with time zone NOT NULL DEFAULT now(),
	admin bool DEFAULT false
);`
	_, err := fs.db.Exec(schema)
	return err
}

// Form returns a Form if its admin field is true
func (fs *FeedbackService) GetAdminForm() (*shuttletracker.Form, error) {
	statement := "SELECT f.id, f.message, f.created, f.admin" +
		" FROM forms f WHERE admin = true;"
	f := &shuttletracker.Form{}
	row := fs.db.QueryRow(statement)
	err := row.Scan(&f.ID, &f.Message, &f.Created, &f.Admin)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrFormNotFound
	}
	return f, err
}

// Form returns a Form by its ID.
func (fs *FeedbackService) GetForm(id int64) (*shuttletracker.Form, error) {
	statement := "SELECT f.message, f.created, f.admin" +
		" FROM forms f WHERE id = $1;"
	f := &shuttletracker.Form{ ID: id }
	row := fs.db.QueryRow(statement, id)
	err := row.Scan(&f.Message, &f.Created, &f.Admin)
	if err == sql.ErrNoRows {
		return nil, shuttletracker.ErrFormNotFound
	}
	return f, err
}

// Forms returns all Forms.
func (fs *FeedbackService) GetForms() ([]*shuttletracker.Form, error) {
	forms := []*shuttletracker.Form{}
	query := "SELECT f.prompt, f.id, f.message, f.created, f.admin" +
		" FROM forms f GROUP BY f.prompt;"
	rows, err := fs.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		f := &shuttletracker.Form{}
		err := rows.Scan(&f.ID, &f.Message, &f.Created, &f.Admin)
		if err != nil {
			return nil, err
		}
		forms = append(forms, f)
	}
	return forms, nil
}

// CreateForm creates a Form.
func (fs *FeedbackService) CreateForm(form *shuttletracker.Form) error {
	if form.Admin == true {
		result, err := fs.db.Exec("DELETE FROM forms WHERE admin = true;")
		if err != nil {
			return err
		}
		n, err := result.RowsAffected()
		if err != nil {
			return err
		}
		log.Debugf(strconv.FormatInt(n, 10) + " stale admin feedback message(s) were successfully deleted")
	}
	statement := "INSERT INTO forms (message, created, admin) VALUES" +
		" ($1, now(), $2) RETURNING id, message, created;"
	row := fs.db.QueryRow(statement, form.Message, form.Admin)
	return row.Scan(&form.ID, &form.Message, &form.Created)
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