package predictor
import schedule
// Stores a simple interface to access table of arrival times stored in database
type TablePredictor struct{
  Table *mgo.Collection
}

// query timetable to get the next N arrival time by select time by stopid with end date > current date > start date, time > current time, weekday = current weekday; 
// limit: NextN should be less than 5 and can go beyond only one day ( check if the result is less than the requested result length)
func (this TablePredictor) getNextN(StopID []string, CurrentTime time.Time , NextN int) []schedule.ArrivalTime{

}

type TimeTable struct{
  ID string `json:"id"`
  StopID string `json:"stopid"`
  StartDate time.Time `json:"startdate"`
  EndDate time.Time `json:"enddate"`
  Time time.Time `json:"time"`
  WeekDay int `json:"weekday"`
}

