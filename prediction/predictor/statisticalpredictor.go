package predictor
import schedule
// Table Predictor stores a simple interface to access the table of arrival times stored in database
type StatisticalPredictor struct{
	History *mgo.Collection // historical data provided
	Velocity *mgo.Collection // velocity data for each segment of a route

}

// limit N <= 5
func (this StatisticalPredictor) getNextN(StopID []string, CurrentTime time.Time , NextN int) []schedule.ArrivalTime{

}

// query when will the shuttle arrive at the next N stops
func (this StatisticalPredictor) getNextN(ShuttleID []string, CurrentLocation Coord, CurrentTime time.Time, NextN int) []schedule.ShuttleArrivalTime{

} 

// if no time within threshold is found, then use the closest one
type RouteVelocity struct{
	ID string `json:"id"`
	CoordID string `json:"coordID"`
	Velocity Coord `json:"velocity"`
	StartTime Time.time `json:"starttime"`
	EndTime Time.time `json:"endtime"`
}