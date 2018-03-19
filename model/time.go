package model

import (
	"strconv"
	"strings"
	"time"
)

type Time struct {
	time time.Time
	Day  int
}

func leftPad(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func (t1 *Time) FromTime(t time.Time) {
	t1.time = t

}

func (t1 *Time) GetTimeString() string {
	s := ""
	hrs := strconv.Itoa(t1.time.Hour())
	mins := strconv.Itoa(t1.time.Minute())
	seconds := strconv.Itoa(t1.time.Second())
	s += leftPad(hrs, "0", 2)
	s += ":"
	s += leftPad(mins, "0", 2)
	s += ":"
	s += leftPad(seconds, "0", 2)
	return s
}

func (t1 Time) After(t2 Time) bool {
	if t1.Day > t2.Day {
		return true
	} else if t1.Day == t2.Day {
		if t1.time.Hour() > t2.time.Hour() {
			return true
		} else if t1.time.Hour() == t2.time.Hour() {
			if t1.time.Minute() > t2.time.Minute() {
				return true
			} else if t1.time.Minute() == t2.time.Minute() {
				if t1.time.Second() > t2.time.Second() {
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
