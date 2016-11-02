package demo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// a demo loading and computing the arrivial time based on time schedule table for each stop
// This is a simple demo, with non structually sound result

// ExampleScheduleFile *
type ExampleScheduleFile struct {
	WeekdayEast  []ExampleSchedule `json:"Weekday_East"`
	WeekendEast  []ExampleSchedule `json:"Weekend_East"`
	WeekdayWest  []ExampleSchedule `json:"Weekday_West"`
	SaturdayWest []ExampleSchedule `json:"Saturday_West"`
	SundayWest   []ExampleSchedule `json:"Sunday_West"`
}

// ExampleSchedule you request this
type ExampleSchedule struct {
	Location string   `json:"location"`
	Times    []string `json:"times"`
}

// Scheduler *
type Scheduler struct {
	Schedule    ExampleScheduleFile `json:"schedule"`
	ArrivalTime ExampleScheduleFile `json:"arrival"`
}

// ClockTime should be 24-hour format
type ClockTime struct {
	Minute int `json:"minute"`
	Hour   int `json:"hour"`
}

//StoClockTime : parse a string with format "HH:SS AM/PM" to ClockTime struct
func StoClockTime(time string) ClockTime {
	times := strings.Split(time, " ")
	times2 := strings.Split(times[0], ":")
	hour, err := strconv.Atoi(times2[0])
	minute, err := strconv.Atoi(times2[1])
	if err != nil {
		log.Fatal("Not Good formated string")
	}
	if times[1] == "PM" {
		hour += 12
	} else if times[1] != "AM" {
		// error
		log.Fatal("Not Good formated string")
	}
	return ClockTime{minute, hour}
}

func ClockTimetoS(c ClockTime) string {
	ampm := "AM"
	if c.Hour > 12 {
		c.Hour -= 12
		ampm = "PM"
	}
	return fmt.Sprintf("%[1]d:%[2]d %[3]s", c.Hour, c.Minute, ampm)
}

func ClockTimeToDuration(c ClockTime) string {
	return fmt.Sprintf("%[1]d hr %[2]d min", c.Hour, c.Minute)
}

// CurrentClockTime *
func CurrentClockTime() ClockTime {
	return StoClockTime(time.Now().Format("03:04 AM"))
}

func ItoClockTime(a int) ClockTime {
	hour := a / 60
	minute := a % 60
	return ClockTime{minute, hour}
}

func ClockTimetoI(a ClockTime) int {
	return a.Hour*60 + a.Minute
}

// ClockTimeDuration calculate a - b
func ClockTimeDuration(a ClockTime, b ClockTime) ClockTime {
	a_I := ClockTimetoI(a)
	b_I := ClockTimetoI(b)
	duration := a_I - b_I // if a < b then we assume that there's one day difference
	if duration < 0 {
		return ItoClockTime(1440 - b_I + a_I)
	}
	return ItoClockTime(duration)
}

// ClockTimeCompare: 0 is a == b, > 0: a > b, < 0: a < b
func ClockTimeCompare(a ClockTime, b ClockTime) int {
	return ClockTimetoI(a) - ClockTimetoI(b)
}

func GetNextNClockTimeString(a ClockTime, b []string, N int) []string {
	if N > len(b) {
		N = len(b)
	}
	index := 0
	for i, value := range b {
		if ClockTimeCompare(a, StoClockTime(value)) < 0 {
			index = i
			break
		}
	}
	result := []string{}
	for i := index; i < index+N; i++ {
		//result = append(result, b[i%len(b)]) // calculate the time
		result = append(result, ClockTimeToDuration(ClockTimeDuration(StoClockTime(b[i%len(b)]), a))) // calculate the duration
	}
	return result
}

// Init s
func Init() Scheduler {
	// load json file:
	raw, err := ioutil.ReadFile("demo/schedule.json")
	if err != nil {
		return Scheduler{}
	}
	schedule := ExampleScheduleFile{}
	json.Unmarshal(raw, &schedule)
	scheduler := Scheduler{}
	scheduler.Schedule = schedule
	return scheduler
}

func (s *Scheduler) CalculateTime() {
	// use current time to calculate the next 3 shuttle arrival time for each location
	s.ArrivalTime = ExampleScheduleFile{}
	// get current clock time
	current_time := CurrentClockTime()
	// let me write some bad code:
	for _, scheduled := range s.Schedule.WeekdayEast {
		arrival_time := ExampleSchedule{}
		arrival_time.Location = scheduled.Location
		arrival_time.Times = GetNextNClockTimeString(current_time, scheduled.Times, 3)
		s.ArrivalTime.WeekdayEast = append(s.ArrivalTime.WeekdayEast, arrival_time)
	}
	for _, scheduled := range s.Schedule.WeekdayWest {
		arrival_time := ExampleSchedule{}
		arrival_time.Location = scheduled.Location
		arrival_time.Times = GetNextNClockTimeString(current_time, scheduled.Times, 3)
		s.ArrivalTime.WeekdayWest = append(s.ArrivalTime.WeekdayWest, arrival_time)
	}
	for _, scheduled := range s.Schedule.WeekendEast {
		arrival_time := ExampleSchedule{}
		arrival_time.Location = scheduled.Location
		arrival_time.Times = GetNextNClockTimeString(current_time, scheduled.Times, 3)
		s.ArrivalTime.WeekendEast = append(s.ArrivalTime.WeekendEast, arrival_time)
	}
	for _, scheduled := range s.Schedule.SaturdayWest {
		arrival_time := ExampleSchedule{}
		arrival_time.Location = scheduled.Location
		arrival_time.Times = GetNextNClockTimeString(current_time, scheduled.Times, 3)
		s.ArrivalTime.SaturdayWest = append(s.ArrivalTime.SaturdayWest, arrival_time)
	}
	for _, scheduled := range s.Schedule.SundayWest {
		arrival_time := ExampleSchedule{}
		arrival_time.Location = scheduled.Location
		arrival_time.Times = GetNextNClockTimeString(current_time, scheduled.Times, 3)
		s.ArrivalTime.SundayWest = append(s.ArrivalTime.SundayWest, arrival_time)
	}
	// so ugly, can't stand it, uhhh
}

type User struct {
	ID                 string            `json:"ID"`
	NotificationMethod map[string]string `json:"noti"`
	SubscribedLocation []string          `json:"loc"`
}

type UserCollection struct {
	Users   map[string]User `json:"user"`
	Content *Scheduler      `json:"scheduler"`
}

func (s *UserCollection) QueryUserInformation(id string) User {
	return s.Users[id]
}

func (s *UserCollection) SubscriberHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user")
	WriteJSON(w, s.MakeANotification(userID))
}

// ArrivalHandler *
func (s *Scheduler) ArrivalHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, s.Schedule)
}

func (s *Scheduler) ArrivalTimeHandler(w http.ResponseWriter, r *http.Request) {
	s.CalculateTime()
	WriteJSON(w, s.ArrivalTime)
}

func UserCollectionInit(scheduler *Scheduler) UserCollection {
	user := User{
		"001",
		map[string]string{"SMS": "5182530000", "EMAIL": "S.A.Jackson.001@ruler.rpi.edu"},
		[]string{"Sage", "Union"}}
	collection := UserCollection{}
	collection.Content = scheduler
	collection.Users = make(map[string]User)
	collection.Users[user.ID] = user
	return collection
}

type NotificationCenter struct {
}

type SMSNotification struct {
	From string `json:"from"`
	To   string `json:"To"`
	Body string `json:"Body"`
}

func Notify(key, value, information string) interface{} {
	var notification interface{}
	switch key {
	case "SMS":
		notification = SMSNotification{"5182530001", value, information}
		return notification
	}
	return notification
}
func (s *Scheduler) GetArrivalTime(locations []string) []ExampleSchedule {
	s.CalculateTime()
	// get current date
	weekday := time.Now().Weekday().String()
	var searchTable []ExampleSchedule
	if weekday == "Sunday" {
		searchTable = append(s.ArrivalTime.SundayWest, s.ArrivalTime.WeekendEast...)
	} else if weekday == "Saturday" {
		searchTable = append(s.ArrivalTime.SaturdayWest, s.ArrivalTime.WeekendEast...)
	} else {
		searchTable = append(s.ArrivalTime.WeekdayEast, s.ArrivalTime.WeekdayWest...)
	}
	// hard coding is not good
	result := []ExampleSchedule{}
	for _, loc := range locations {
		schedule := ExampleSchedule{}
		breakFlag := false
		for _, loc2 := range searchTable {
			if loc2.Location == loc {
				for _, t := range loc2.Times {
					if len(schedule.Times) >= 3 {
						breakFlag = true
						break
					}
					schedule.Times = append(schedule.Times, t)
				}

			}
			if breakFlag {
				break
			}
		}
		schedule.Location = loc
		result = append(result, schedule)
	}
	return result
}

func GetReadable(schedule ExampleSchedule) string {
	var time_string string
	for i, time := range schedule.Times {
		if i != 0 {
			time_string += ","
		}
		time_string += time
	}
	return schedule.Location + " will arrive in: " + time_string
}
func GetReadableInformation(key string, schedules []ExampleSchedule) string {
	var result string
	switch key {
	case "SMS":
		for _, s := range schedules {
			result += GetReadable(s) + ","
		}
	}
	return result
}
func (s *UserCollection) MakeANotification(userID string) interface{} {
	var result []interface{}
	for key, value := range s.Users[userID].NotificationMethod {
		information := GetReadableInformation(key, s.Content.GetArrivalTime(s.Users[userID].SubscribedLocation))
		result = append(result, Notify(key, value, information))
	}
	return result
}

// WriteJSON writes the data as JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}

func main() {
	s := Init()
	s.CalculateTime()
}
