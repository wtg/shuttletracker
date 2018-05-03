package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/wtg/shuttletracker"
)

const url = "postgres://postgres@localhost/shuttletracker_test?sslmode=disable"

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
	_, err = db.Exec("drop schema public cascade; create schema public authorization postgres;")
	if err != nil {
		t.Fatalf("unable to clear database: %s", err)
	}
}

// nolint: gocyclo
func TestCreateLocation(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	// insert test data
	vehicle := &shuttletracker.Vehicle{
		Name:      "test vehicle",
		Enabled:   false,
		TrackerID: "tracker1",
	}
	err := pg.CreateVehicle(vehicle)
	if err != nil {
		t.Fatalf("unable to create Vehicle: %s", err)
	}

	location := &shuttletracker.Location{
		TrackerID: "tracker1",
		Latitude:  1.1,
		Longitude: 1.2,
		Heading:   1.3,
		Speed:     1.4,
		RouteID:   nil,
		Time:      time.Now(),
	}
	err = pg.CreateLocation(location)
	if err != nil {
		t.Fatalf("unable to create Location: %s", err)
	}

	// retrieve test data
	actual, err := pg.LatestLocation(vehicle.ID)
	if err != nil {
		t.Fatalf("unable to get latest Location: %s", err)
	}

	if location.ID != actual.ID {
		t.Errorf("got ID %d, expected %d", actual.ID, location.ID)
	}
	if location.TrackerID != actual.TrackerID {
		t.Errorf("got tracker ID %s, expected %s", actual.TrackerID, location.TrackerID)
	}
	if location.Latitude-actual.Latitude > 0.0000001 {
		t.Errorf("got latitude %f, expected %f", actual.Latitude, location.Latitude)
	}
	if location.Longitude-actual.Longitude > 0.0000001 {
		t.Errorf("got longitude %f, expected %f", actual.Longitude, location.Longitude)
	}
	if location.Heading-actual.Heading > 0.0000001 {
		t.Errorf("got heading %f, expected %f", actual.Heading, location.Heading)
	}
	if location.Speed-actual.Speed > 0.0000001 {
		t.Errorf("got speed %f, expected %f", actual.Speed, location.Speed)
	}
	if location.RouteID != actual.RouteID {
		t.Errorf("got route ID %d, expected %d", actual.RouteID, location.RouteID)
	}
	if location.Time.Sub(actual.Time).Nanoseconds() > 1000 {
		t.Errorf("got time %v, expected %v", actual.Time, location.Time)
	}
	if time.Since(actual.Time).Seconds() > 1 {
		t.Errorf("got created %v, which is too old", actual.Created)
	}
}

// nolint: gocyclo
func TestDeleteLocationsBefore(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	// insert test data
	vehicle := &shuttletracker.Vehicle{
		Name:      "test vehicle",
		Enabled:   false,
		TrackerID: "tracker1",
	}
	err := pg.CreateVehicle(vehicle)
	if err != nil {
		t.Fatalf("unable to create Vehicle: %s", err)
	}

	location1 := &shuttletracker.Location{
		TrackerID: "tracker1",
		Latitude:  1.1,
		Longitude: 1.2,
		Heading:   1.3,
		Speed:     1.4,
		RouteID:   nil,
		Time:      time.Now(),
	}
	err = pg.CreateLocation(location1)
	if err != nil {
		t.Fatalf("unable to create Location: %s", err)
	}
	location2 := &shuttletracker.Location{
		TrackerID: "tracker1",
		Latitude:  1.1,
		Longitude: 1.2,
		Heading:   1.3,
		Speed:     1.4,
		RouteID:   nil,
		Time:      time.Now(),
	}
	err = pg.CreateLocation(location2)
	if err != nil {
		t.Fatalf("unable to create Location: %s", err)
	}
	lastTime := time.Now()
	location3 := &shuttletracker.Location{
		TrackerID: "tracker1",
		Latitude:  1.1,
		Longitude: 1.2,
		Heading:   1.3,
		Speed:     1.4,
		RouteID:   nil,
		Time:      lastTime,
	}
	err = pg.CreateLocation(location3)
	if err != nil {
		t.Fatalf("unable to create Location: %s", err)
	}

	// retrieve test data
	actuals, err := pg.LocationsSince(vehicle.ID, time.Time{})
	if err != nil {
		t.Fatalf("unable to get Locations: %s", err)
	}

	if len(actuals) != 3 {
		t.Fatalf("got %d Locations, expected 3", len(actuals))
	}

	// delete Locations
	n, err := pg.DeleteLocationsBefore(lastTime)
	if err != nil {
		t.Fatalf("unable to delete Locations: %s", err)
	}
	if n != 2 {
		t.Errorf("deleted %d Locations, expected 2", n)
	}

	// retrieve test data
	actuals, err = pg.LocationsSince(vehicle.ID, time.Time{})
	if err != nil {
		t.Fatalf("unable to get Locations: %s", err)
	}

	if len(actuals) != 1 {
		t.Fatalf("got %d Locations, expected 1", len(actuals))
	}
}
