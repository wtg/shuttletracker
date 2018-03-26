package database

import (
	"github.com/wtg/shuttletracker/model"
	"testing"

	"time"
)

/*
  Note: A couple of these tests fail around midnight, specifically the time after function
  has several issues due to the time changing from 24:00 to 00:00 this does not affect the
  functionality of the shuttle tracker
*/

func TestSetRouteActiveStatusSameDay(t *testing.T) {
	//Test when the day is the same, and time varies and should be active
	interval := []model.Time{}
	t1 := model.Time{
		Day:   time.Now().Weekday(),
		Time:  time.Now().Add(-5 * time.Second),
		State: 1,
	}
	t2 := model.Time{
		Day:   time.Now().Weekday(),
		Time:  time.Now().Add(1 * time.Hour),
		State: 0,
	}
	interval = append(interval, t1, t2)

	TestRoute := model.Route{
		ID:           "Test",
		Name:         "Test Route",
		TimeInterval: interval,
	}
  SetRouteActiveStatus(&TestRoute)
	if TestRoute.Active {
		t.Errorf("Route should be active but is not %+v", TestRoute)
	}
}

func TestSetRouteActiveStatusSimilarTime(t *testing.T) {
	//Test when the day is the same, and time varies and should be active
	interval := []model.Time{}
	t1 := model.Time{
		Day:   time.Now().Weekday(),
		Time:  time.Now().Add(-5 * time.Second),
		State: 1,
	}
	t2 := model.Time{
		Day:   time.Now().Weekday(),
		Time:  time.Now().Add(5 * time.Second),
		State: 0,
	}
	interval = append(interval, t1, t2)

	TestRoute := model.Route{
		ID:           "Test",
		Name:         "Test Route",
		TimeInterval: interval,
	}
	SetRouteActiveStatus(&TestRoute)
	if !TestRoute.Active {
		t.Errorf("Route should be active but is not")
	}
}

//Creates some test routes and tests if they should be active
func TestSetRouteActiveStatusDiffDay(t *testing.T) {
	//Test when the day is explicitly different and should be active
	interval := []model.Time{}
	t1 := model.Time{
		Day:   time.Now().Weekday() - 1,
		Time:  time.Now().Add(-5 * time.Minute),
		State: 1,
	}
	t2 := model.Time{
		Day:   time.Now().Weekday(),
		Time:  time.Now().Add(5 * time.Minute),
		State: 0,
	}
	interval = append(interval, t1, t2)

	TestRoute := model.Route{
		ID:           "Test",
		Name:         "Test Route",
		TimeInterval: interval,
	}

	SetRouteActiveStatus(&TestRoute)
	if !TestRoute.Active {
		t.Errorf("Route should be active but is not %+v", TestRoute)
	}

	//Test when the day is explicitly different and should not be active
	interval = []model.Time{}
	t1 = model.Time{
		Day:   time.Now().Weekday() - 1,
		Time:  time.Now().Add(-5 * time.Minute),
		State: 0,
	}
	t2 = model.Time{
		Day:   time.Now().Weekday(),
		Time:  time.Now().Add(5 * time.Minute),
		State: 1,
	}
	interval = append(interval, t1, t2)
	TestRoute.TimeInterval = interval
	SetRouteActiveStatus(&TestRoute)
	if TestRoute.Active {
		t.Errorf("Route should not be active but is %+v", TestRoute)
	}
}
