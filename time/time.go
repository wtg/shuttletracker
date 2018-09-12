package time

import (
	"fmt"
	"sort"
	"time"
)

// Time struct is used to store route enable/disable times, so we can compare times without a date
type Time struct {
	Time  time.Time    `json:"time"`
	Day   time.Weekday `json:"day"`
	State int          `json:"on"`
}

// FromTime is used to create a new model.Time object from a normal gotime object
func (t1 *Time) FromTime(t time.Time) {
	t1.Time = t
}

// GetTimeString returns the string version of the time represented by the struct
func (t1 *Time) GetTimeString() string {
	s := ""
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

// byTime is an interface used to sort times.
type byTime []Time

func (a byTime) Len() int {
	return len(a)
}

func (a byTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byTime) Less(i, j int) bool {
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

// Sort sorts a slice of Times.
// We cannot use sort.Slice because it does not exist in Go 1.7.
func Sort(times []Time) {
	sort.Sort(byTime(times))
}
