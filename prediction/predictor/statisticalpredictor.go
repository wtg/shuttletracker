package predictor
import schedule
// Stores a simple interface to access table of arrival times stored in database
type StatisticalPredictor struct{
  History *mgo.Collection // historical data provided
  Velocity *mgo.Collection // velocity data for each segment of a route
}

func (this StatisticalPredictor) getNextN(StopID []string, CurrentTime time.Time , NextN int) []schedule.ArrivalTime{
  // limit N <= 5
}

// Get prediction for shuttle arrive at the next N stops
func (this StatisticalPredictor) getNextN(ShuttleID []string, CurrentLocation Coord, CurrentTime time.Time, NextN int) []schedule.ShuttleArrivalTime{

} 

type RouteVelocity struct{
  ID string `json:"id"`
  CoordID string `json:"coordID"`
  Velocity Coord `json:"velocity"`
  StartTime Time.time `json:"starttime"`
  EndTime Time.time `json:"endtime"`
  // if no time within threshold is found, then use the closest one
}
