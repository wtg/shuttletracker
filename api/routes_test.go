package api

import(   "testing"
  "github.com/wtg/shuttletracker/model"
  "github.com/wtg/shuttletracker/api"
  "fmt"

  "time"
)

func TestRouteIsActive(t *testing.T) {

  api := api.API{
  }

  testInterval := []model.WeekTime{}
  t1 := model.WeekTime{
    Day: time.Now().Weekday(),
    Time: time.Now().Add(-1*time.Minute),
    State: 0,
  }
  testInterval = append(testInterval, t1)
  testRoute := model.Route{
    ID: "test",
    Name: "testRoute",
    TimeInterval: testInterval,
  }
  
  fmt.Printf("test: ", api.RouteIsActive(&testRoute))



}
