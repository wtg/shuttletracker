// Utilities for testing the postgres package.
package postgres

import (
	"database/sql"
	"testing"
)

const url = "postgres://localhost/shuttletracker_test?sslmode=disable"

func setUpPostgres(t *testing.T) *Postgres {
	// determine if this database is clean
	db, err := sql.Open("postgres", url)
	if err != nil {
		t.Fatalf("unable to open database: %s", err)
	}
	row := db.QueryRow("select count(*) from information_schema.tables where table_schema = 'public';")
	var numTables int
	err = row.Scan(&numTables)
	if err != nil {
		t.Fatalf("unable to scan: %s", err)
	}
	if numTables != 0 {
		t.Fatalf("database is not empty")
	}

	pg, err := New(Config{URL: url})
	if err != nil {
		t.Fatalf("unable to create Postgres: %s", err)
	}

	return pg
}

func tearDownPostgres(t *testing.T) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		t.Fatalf("unable to open database: %s", err)
	}
	_, err = db.Exec("drop schema public cascade; create schema public;")
	if err != nil {
		t.Fatalf("unable to clear database: %s", err)
	}
}
