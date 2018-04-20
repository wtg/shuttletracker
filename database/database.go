package database

import (
	"errors"
	"time"

	"github.com/wtg/shuttletracker/model"
)

//SetRouteActiveStatus determines if a given route is active based on its schedule intervals and the time given, then updates the object in the parameter
func SetRouteActiveStatus(r *model.Route, t time.Time) {

	//This is a time offset, to ensure routes are activated on the minute they are assigned activate
	var currentTime model.Time
	currentTime.FromTime(t)
	currentTime.Day = time.Now().Weekday()
	state := -1

	if r.TimeInterval == nil || len(r.TimeInterval) == 1 {
		state = 1
	}
	for idx, val := range r.TimeInterval {
		//If it is the last in the time list (latest time for the week) use this index
		if idx >= len(r.TimeInterval)-1 {
			state = val.State
			break
		} else {
			if currentTime.After(val) && currentTime.After(r.TimeInterval[idx+1]) {
				continue
			}
			state = val.State
			break
		}
	}

	r.Active = (state == 1 || state == -1)

}

// Database is an interface that can be implemented by a database backend.
type Database interface {
	// Stops
	CreateStop(stop *model.Stop) error
	DeleteStop(stopID string) error
	GetStops() ([]model.Stop, error)
	// GetStopsForRoute(routeID string) ([]model.Stop, error)
	// ModifyStop(stop *model.Stop) error

}

var (
	// ErrUpdateNotFound indicates that an Update is not in the database.
	ErrUpdateNotFound = errors.New("Update not found")
)
