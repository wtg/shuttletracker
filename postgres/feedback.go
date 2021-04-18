package postgres

import (
	"database/sql"
	"strconv"
	"fmt"
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
	admin_shuttletracker bool DEFAULT false
);`
	_, err := fs.db.Exec(schema)
	return err
}

// Form returns a Form if its admin field is true
func (fs *FeedbackService) GetAdminForm() (*shuttletracker.Form) {
	statement := "SELECT f.id, f.message, f.created, f.admin_shuttletracker" +
		" FROM forms f WHERE admin_shuttletracker = true;"
	f := &shuttletracker.Form{}
	row := fs.db.QueryRow(statement)
	row.Scan(&f.ID, &f.Message, &f.Created, &f.Admin)
	return f
}

// Form returns a Form by its ID.
func (fs *FeedbackService) GetForm(id int64) (*shuttletracker.Form, error) {
	statement := "SELECT f.message, f.created, f.admin_shuttletracker" +
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
	// Add GROUP BY f.prompt to sort the table with
	query := "SELECT f.prompt, f.id, f.message, f.created, f.admin_shuttletracker" +
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
		fmt.Println("admin is true!!!!!")
		result, err := fs.db.Exec("DELETE FROM forms WHERE admin_shuttletracker = true;")
		fmt.Println("111111!")
		if err != nil {
			fmt.Println("eeeeeee!")
			fmt.Println(err)
			return err
		}
		fmt.Println("121!")
		n, err := result.RowsAffected()
		fmt.Println("1113!")
		if err != nil {
			return err
		}
		log.Debugf(strconv.FormatInt(n, 10) + " stale admin feedback message(s) were successfully deleted")
	}
	// Add prompt when inserting value to the table
	statement := "INSERT INTO forms (prompt, message, created, admin_shuttletracker) VALUES" +
		" ($1, $2, now(), $3) RETURNING id, message, created;"
	// Use GetAdminForm function to get the current prompt visible to user
	fmt.Println("111111!")
	msg := fs.GetAdminForm()
	fmt.Println("22222!")
	row := fs.db.QueryRow(statement, msg.Message, form.Message, form.Admin)
	fmt.Println("33333!")
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