package postgres

import (
	"database/sql"
	"testing"
)

func TestPostgres(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	cfg := Config{
		URL: "postgres://postgres@localhost/shuttletracker_test?sslmode=disable",
	}
	// determine if this database is clean
	db, err := sql.Open("postgres", cfg.URL)
	row := db.QueryRow("select count(*) from information_schema.tables where table_schema = 'public';")
	var numTables int
	err = row.Scan(&numTables)
	if err != nil {
		t.Fatalf("unable to scan: %s", err)
	}

	if numTables != 0 {
		t.Fatalf("database is not empty")
	}

	// p, err := New(cfg)
	// if err != nil {
	// 	t.Fatalf("unexpected error: %s", err)
	// }

}
