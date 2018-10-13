package postgres

import (
	"testing"
	"time"

	"github.com/wtg/shuttletracker"
)

func TestCreateEmptySchedule(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	route := &shuttletracker.Route{
		Name:     "Test Route",
		Schedule: shuttletracker.RouteSchedule{},
	}
	err := pg.CreateRoute(route)
	if err != nil {
		t.Fatalf("unable to create Route: %s", err)
	}

	// make sure the schedule is empty and route is active
	if len(route.Schedule) != 0 {
		t.Error("schedule is not empty")
	}
	if !route.Active {
		t.Error("route is not active")
	}

	// retrieve the route again and see if it's empty
	route, err = pg.Route(route.ID)
	if err != nil {
		t.Fatalf("unable to get Route: %s", err)
	}
	if len(route.Schedule) != 0 {
		t.Error("schedule is not empty")
	}
	if !route.Active {
		t.Error("route is not active")
	}
}

// nolint: gocyclo
func TestCreateSchedule(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	route := &shuttletracker.Route{
		Name: "Test Route",
		Schedule: shuttletracker.RouteSchedule{
			shuttletracker.RouteActiveInterval{
				StartDay:  time.Sunday,
				StartTime: time.Date(0, 1, 0, 0, 0, 0, 0, time.UTC),
				EndDay:    time.Saturday,
				EndTime:   time.Date(0, 1, 0, 23, 59, 59, 0, time.UTC),
			},
		},
	}
	err := pg.CreateRoute(route)
	if err != nil {
		t.Fatalf("unable to create Route: %s", err)
	}

	// make sure the schedule is the correct size and route is active
	if len(route.Schedule) != 1 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}

	// retrieve the route and check again
	route, err = pg.Route(route.ID)
	if err != nil {
		t.Fatalf("unable to get Route: %s", err)
	}
	if len(route.Schedule) != 1 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}

	// check that we can get it from all routes
	routes, err := pg.Routes()
	if err != nil {
		t.Fatalf("unable to get Routes: %s", err)
	}
	if len(routes) != 1 {
		t.Fatalf("wrong Routes length: %d", len(routes))
	}
	route = routes[0]
	if len(route.Schedule) != 1 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}
}

// nolint: gocyclo
func TestModifySchedule(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	route := &shuttletracker.Route{
		Name:     "Test Route",
		Schedule: shuttletracker.RouteSchedule{},
	}
	err := pg.CreateRoute(route)
	if err != nil {
		t.Fatalf("unable to create Route: %s", err)
	}

	// make sure the schedule is the correct size and route is active
	if len(route.Schedule) != 0 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}

	// modify the route schedule and check again
	route.Schedule = shuttletracker.RouteSchedule{
		shuttletracker.RouteActiveInterval{
			StartDay:  time.Sunday,
			StartTime: time.Date(0, 1, 0, 0, 0, 0, 0, time.UTC),
			EndDay:    time.Saturday,
			EndTime:   time.Date(0, 1, 0, 23, 59, 59, 0, time.UTC),
		},
	}
	err = pg.ModifyRoute(route)
	if err != nil {
		t.Fatalf("unable to modify Route: %s", err)
	}
	if len(route.Schedule) != 1 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}

	route, err = pg.Route(route.ID)
	if err != nil {
		t.Fatalf("unable to get Route: %s", err)
	}
	if len(route.Schedule) != 1 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}

	// modify the route schedule and check again
	route.Schedule = shuttletracker.RouteSchedule{
		shuttletracker.RouteActiveInterval{
			StartDay:  time.Sunday,
			StartTime: time.Date(0, 1, 0, 0, 0, 0, 0, time.UTC),
			EndDay:    time.Tuesday,
			EndTime:   time.Date(0, 1, 0, 23, 59, 59, 0, time.UTC),
		},
		shuttletracker.RouteActiveInterval{
			StartDay:  time.Wednesday,
			StartTime: time.Date(0, 1, 0, 0, 0, 0, 0, time.UTC),
			EndDay:    time.Saturday,
			EndTime:   time.Date(0, 1, 0, 23, 59, 59, 0, time.UTC),
		},
	}
	err = pg.ModifyRoute(route)
	if err != nil {
		t.Fatalf("unable to modify Route: %s", err)
	}
	if len(route.Schedule) != 2 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}

	route, err = pg.Route(route.ID)
	if err != nil {
		t.Fatalf("unable to get Route: %s", err)
	}
	if len(route.Schedule) != 2 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}

	// delete the route schedule and check again
	route.Schedule = shuttletracker.RouteSchedule{}
	err = pg.ModifyRoute(route)
	if err != nil {
		t.Fatalf("unable to modify Route: %s", err)
	}
	if len(route.Schedule) != 0 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}

	route, err = pg.Route(route.ID)
	if err != nil {
		t.Fatalf("unable to get Route: %s", err)
	}
	if len(route.Schedule) != 0 {
		t.Errorf("wrong schedule length: %d", len(route.Schedule))
	}
	if !route.Active {
		t.Error("route is not active")
	}
}

// This test probably doesn't work properly around midnight...
// nolint: gocyclo
func TestActiveTransition(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pg := setUpPostgres(t)
	defer tearDownPostgres(t)

	now := time.Now()
	route := &shuttletracker.Route{
		Name: "Test Route",
		Schedule: shuttletracker.RouteSchedule{
			shuttletracker.RouteActiveInterval{
				StartDay:  now.Weekday(),
				StartTime: now,
				EndDay:    now.Weekday(),
				EndTime:   now.Add(time.Second),
			},
			shuttletracker.RouteActiveInterval{
				StartDay:  now.Weekday(),
				StartTime: now.Add(5 * time.Second),
				EndDay:    now.Weekday(),
				EndTime:   now.Add(10 * time.Second),
			},
		},
	}
	err := pg.CreateRoute(route)
	if err != nil {
		t.Fatalf("unable to create Route: %s", err)
	}
	if !route.Active {
		t.Error("route is not active")
	}

	// wait two seconds and then check again...
	time.Sleep(2 * time.Second)
	route, err = pg.Route(route.ID)
	if err != nil {
		t.Fatalf("unable to get Route: %s", err)
	}
	if route.Active {
		t.Error("route is active")
	}

	// wait and then check again...
	time.Sleep(4 * time.Second)
	route, err = pg.Route(route.ID)
	if err != nil {
		t.Fatalf("unable to get Route: %s", err)
	}
	if !route.Active {
		t.Error("route is not active")
	}

	// wait and then check again...
	time.Sleep(5 * time.Second)
	route, err = pg.Route(route.ID)
	if err != nil {
		t.Fatalf("unable to get Route: %s", err)
	}
	if route.Active {
		t.Error("route is active")
	}
}
