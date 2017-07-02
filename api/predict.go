package api

// Arrival Time serving what's the time to next N stops for one shuttle
import (
	"fmt"
	"math"

	"strconv"

	"bytes"

	"github.com/wtg/shuttletracker/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetArrivalTime is experimental
func GetArrivalTime(update *model.VehicleUpdate, routes *mgo.Collection, stops *mgo.Collection) string {
	if i, err := strconv.ParseFloat(update.Speed, 64); i > 5.0 && err == nil {
		route := model.Route{}
		routes.Find(bson.M{"id": "582f2794e05a0b9c1f2948fa"}).One(&route)
		// get closest segment
		x0, err := strconv.ParseFloat(update.Lat, 64)
		if err != nil {
			panic("Parsing Error")
		}
		y0, err := strconv.ParseFloat(update.Lng, 64)
		if err != nil {
			panic("Parsing Error")
		}
		minimumLen := 1000.0 // bad code
		ShuttleSegment := 0
		for i := 0; i < len(route.Duration); i++ {
			x1 := route.Duration[i].Start.Latitude
			y1 := route.Duration[i].Start.Longitude
			x2 := route.Duration[i].End.Latitude
			y2 := route.Duration[i].End.Longitude
			// compute the distance between a point and a line
			length := math.Abs((x2-x1)*(y1-y0)-(x1-x0)*(y2-y1)) / math.Sqrt(math.Pow(x2-x1, 2)+math.Pow(y2-y1, 2))
			if length < minimumLen {
				minimumLen = length
				ShuttleSegment = i
			}
		}
		// when shuttle segment is not found, return N/A
		if ShuttleSegment >= 0 && ShuttleSegment < len(route.Duration) {
			fmt.Printf("ID = %s, Segment = %d (%f, %f), duration = %f\n", update.VehicleID, ShuttleSegment, route.Duration[ShuttleSegment].Start.Latitude, route.Duration[ShuttleSegment].Start.Longitude, route.Duration[ShuttleSegment].Duration)
		}
		allstops := []model.Stop{}
		stops.Find(bson.M{"routeId": route.ID}).All(&allstops)
		var buffer bytes.Buffer
		for _, i := range allstops {
			Timecost := 0.0
			if i.SegmentIndex == ShuttleSegment {
				buffer.WriteString(fmt.Sprintf("Arrived at: %s;", i.Name))
			}
			if i.SegmentIndex > ShuttleSegment { // this is the next
				for index := ShuttleSegment; index < i.SegmentIndex; index++ {
					Timecost += route.Duration[index].Duration
				}
				fmt.Printf("Cost: %f, Next Stop Name: %s, Segment:%d(%f, %f)\n", Timecost, i.Name, i.SegmentIndex, route.Duration[i.SegmentIndex].Start.Latitude, route.Duration[i.SegmentIndex].Start.Longitude)

				buffer.WriteString(fmt.Sprintf("Stop: %s (%f);", i.Name, Timecost))
			}
		}
		return buffer.String()
	}
	return "N/A"
}
