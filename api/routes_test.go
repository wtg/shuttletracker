package api

import (
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/model"
	"testing"

	"time"
)

func TestTimeAfter(t *testing.T) {
	t1 := time.Now().Add(-time.Second)
	t2 := time.Now().Add(time.Second)
	t3 := t1.Add(time.Hour)
	t4 := t1.Add(time.Minute)

	if api.TimeAfter(t1, t2) {
		t.Errorf("t1 should be before t2")
	}
	if !api.TimeAfter(t2, t1) {
		t.Errorf("t2 should be after t1")
	}
	if api.TimeAfter(t2, t2) {
		t.Errorf("t2 should not be after itself")
	}
	if !api.TimeAfter(t3, t1) {
		t.Errorf("t3 should be after t1")
	}
	if !api.TimeAfter(t4, t1) {
		t.Errorf("t4 should be after t1")
	}

}

func TestRouteIsActiveSameDay(t *testing.T) {
	api := api.API{}

	//Test when the day is the same, and time varies and should be active
	interval := []model.WeekTime{}
	t1 := model.WeekTime{
		Day:   time.Now().Weekday(),
		Time:  time.Now().Add(-5 * time.Hour),
		State: 1,
	}
	t2 := model.WeekTime{
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
	if !api.RouteIsActive(&TestRoute) {
		t.Errorf("Route should be active but is not")
	}
}

//Creates some test routes and tests if they should be active
func TestRouteIsActiveDiffDay(t *testing.T) {
	api := api.API{}

	//Test when the day is explicitly different and should be active
	interval := []model.WeekTime{}
	t1 := model.WeekTime{
		Day:   time.Now().Weekday() - 1,
		Time:  time.Now().Add(-5 * time.Minute),
		State: 1,
	}
	t2 := model.WeekTime{
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

	if !api.RouteIsActive(&TestRoute) {
		t.Errorf("Route should be active but is not")
	}

	//Test when the day is explicitly different and should not be active
	interval = []model.WeekTime{}
	t1 = model.WeekTime{
		Day:   time.Now().Weekday() - 1,
		Time:  time.Now().Add(-5 * time.Minute),
		State: 0,
	}
	t2 = model.WeekTime{
		Day:   time.Now().Weekday(),
		Time:  time.Now().Add(5 * time.Minute),
		State: 1,
	}
	interval = append(interval, t1, t2)
	TestRoute.TimeInterval = interval
	if api.RouteIsActive(&TestRoute) {
		t.Errorf("Route should not be active but is")
	}
}
