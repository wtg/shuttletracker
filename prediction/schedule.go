package schedule

import "time"

// Struct for arrival time for next N shuttles
//  stores data for only 1 shuttle
type ShuttleArrivalTime struct{
  ID string `json:"id"`
  ShuttleID string `json:"shuttleid"`
  Created time.Time `json:"created"`
  Arrival []ArrivalTime `json:"arrival"`
}

// Struct for shuttle arrival time for next stop
type ArrivalTime struct {
  ID string `json:"id"`
  StopID string `json:"stopid"`
  Created time.Time `json:"created"`
  Arrival []time.Time `json:"arrival`
}

// Interface for generating/formatting ArrivalTime
type ArrivalPredictor interface {
  getNextN(StopID []string, CurrentTime time.Time, NextN int) []ArrivalTime
}

type ShuttleArrivalPredictor interface{
  getNextN(StopID []string, CurrentTime time.Time, NextN int) []ShuttleArrivalTime
}

