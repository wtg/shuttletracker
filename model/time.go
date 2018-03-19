package model

import (
	"fmt"
	"time"
)

// Time struct is used to store route enable/disable times, so we can compare times without a date
type Time struct {
	Time  time.Time    `json:time`
	Day   time.Weekday `json:day`
	State int          `json:on`
}

// FromTime is used to create a new model.Time object from a normal gotime object
func (t1 *Time) FromTime(t time.Time) {
	t1.Time = t

}

// GetTimeString returns the string version of the time represented by the struct
func (t1 *Time) GetTimeString() string {
	s := ""

	fmt.Printf("|%06d|%6d|\n", 12, 345)
	s += fmt.Sprintf("%02d", t1.Time.Hour())
	s += ":"
	s += fmt.Sprintf("%02d", t1.Time.Minute())
	s += ":"
	s += fmt.Sprintf("%02d", t1.Time.Second())
	return s
}

// After returns true if the time is after the parameter time, false otherwise
func (t1 Time) After(t2 Time) bool {
	if t1.Day > t2.Day {
		return true
	} else if t1.Day == t2.Day {
		if t1.Time.Hour() > t2.Time.Hour() {
			return true
		} else if t1.Time.Hour() == t2.Time.Hour() {
			if t1.Time.Minute() > t2.Time.Minute() {
				return true
			} else if t1.Time.Minute() == t2.Time.Minute() {
				if t1.Time.Second() > t2.Time.Second() {
					return true
				}
			}
		}

	}
	return false
}

//ByTime is an interface used to sort times
type ByTime []Time

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {
	if a[i].Day < a[j].Day {
		return true
	} else if a[i].Day > a[j].Day {
		return false
	} else {
		if a[j].After(a[i]) {
			return true
		}
		return false
	}
}
