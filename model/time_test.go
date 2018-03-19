package model

import (
	"testing"
  "time"
  "sort"
  "github.com/wtg/shuttletracker/model"
)

func TestCreateTime(t *testing.T) {
  timeTest,_ := time.Parse("15:04:05","12:05:06")
  var t1 model.Time
  t1.FromTime(timeTest)
  if !(t1.GetTimeString() == "12:05:06"){
    t.Errorf("Time string did not match anticipated string")
  }
}

func TestTimeAfter(t *testing.T) {
  timeTest,_ := time.Parse("15:04:05","12:05:06")
  timeTestTwo := timeTest.Add(1*time.Minute)
  var t1 model.Time
  t1.FromTime(timeTest)
  var t2 model.Time
  t2.FromTime(timeTestTwo)
  if(t1.After(t2)){
    t.Errorf("t2 should be after t1")
  }
  if(!t2.After(t1)){
    t.Errorf("t1 should not be after t2")
  }
  if(t1.After(t1)){
    t.Errorf("a time cannot be after itself")
  }

}

func TestTimeAfterWithDay(t *testing.T) {
  timeTest,_ := time.Parse("15:04:05","12:05:06")
  timeTestTwo := timeTest.Add(1*time.Minute)
  var t1 model.Time
  t1.FromTime(timeTest)
  var t2 model.Time
  t2.FromTime(timeTestTwo)
  t2.Day = 1
  t1.Day = 2
  if(!t1.After(t2)){
    t.Errorf("t1 should be after t2")
  }
  if(t2.After(t1)){
    t.Errorf("t2 should not be after t1")
  }
  if(t1.After(t1)){
    t.Errorf("a time cannot be after itself")
  }

}

func TestSorting(t *testing.T){
  timeTest,_ := time.Parse("15:04:05","12:05:06")
  timeTestTwo := timeTest.Add(1*time.Minute)
  timeTestThree := timeTest.Add(-1*time.Minute)
  var t1 model.Time
  t1.FromTime(timeTest)
  var t2 model.Time
  t2.FromTime(timeTestTwo)
  var t3 model.Time
  t3.FromTime(timeTestThree)
  t2.Day = 1
  t1.Day = 2
  t3.Day = 1
  var times []model.Time
  times = append(times, t1,t2,t3)

  sort.Sort(model.ByTime(times))
  for idx,it := range times{
    if(idx < len(times) - 1){
      if !times[idx+1].After(it){
        t.Errorf("Times %d and %d not in sorted order", idx, idx+1)
      }
    }
  }
}
