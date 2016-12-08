package schedule

import "time"

// Arrival Time serving what's the time to next N stops for one shuttle
type ShuttleArrivalTime struct{
	ID string `json:"id"`
	ShuttleID string `json:"shuttleid"`
	Created time.Time `json:"created"`
	Arrival []ArrivalTime `json:"arrival"` // this stores only the arrival time for stops for this specific shuttle
}

// Arrival Time serving what's the time to next N shuttle arrival for one stop
type ArrivalTime struct {
	ID string `json:"id"`
	StopID string `json:"stopid"`
	Created time.Time `json:"created"`
	Arrival []time.Time `json:"arrival`
}

// an interface for generating ArrivalTime
type ArrivalPredictor interface {
	getNextN(StopID []string, CurrentTime time.Time, NextN int) []ArrivalTime
}

type ShuttleArrivalPredictor interface{
	getNextN(StopID []string, CurrentTime time.Time, NextN int) []ShuttleArrivalTime
}

