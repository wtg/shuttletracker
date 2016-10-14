package schedule

type Schedule struct {
	ID string `json:"id"`
	StopID string `json:"stopid"`
	Times []Times.time `json:"times"` 
}

type ScheduleMap struct{
	Schedules map[string]Schedule // key: stopid, value: schedule
}

type ArrivalTime struct{
	TimeStamp string
	Arrival map[string][]int // each stop contains one 
}

