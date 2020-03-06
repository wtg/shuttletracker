package smooth

import (
	"math"
	"time"

	"github.com/wtg/shuttletracker"
)

type SmoothTracking struct {
	ms          shuttletracker.ModelService

type Prediction struct {
	VehicleID int64
	Point     shuttletracker.Point
	Index     int
}

const earthRadius = 6371000.0 // meters

func haversine(theta float64) float64 {
	return (1 - math.Cos(theta)) / 2
}

func toRadians(n float64) float64 {
	return n * math.Pi / 180
}

func DistanceBetween(p1, p2 shuttletracker.Point) float64 {
	lat1Rad := toRadians(p1.Latitude)
	lon1Rad := toRadians(p1.Longitude)
	lat2Rad := toRadians(p2.Latitude)
	lon2Rad := toRadians(p2.Longitude)

	return 2 * earthRadius * math.Asin(math.Sqrt(
		haversine(lat2Rad-lat1Rad)+
			math.Cos(lat1Rad)*math.Cos(lat2Rad)*
				haversine(lon2Rad-lon1Rad)))
}

// Returns the index of the closest point on the route to the given latitude and longitude coordinates
func ClosestPointTo(latitude, longitude float64, route *shuttletracker.Route) int {
	index := 0
	minDistance := math.Inf(1)
	for i, point := range route.Points {
		distance := DistanceBetween(point, shuttletracker.Point{Latitude: latitude, Longitude: longitude})
		if distance < minDistance {
			minDistance = distance
			index = i
		}
	}
	return index
}

// Naive algorithm to predict the position a shuttle is at, given the last update received
// Returns the index of the point the shuttle would be at on its route
func NaivePredictPosition(vehicle *shuttletracker.Vehicle, lastUpdate *shuttletracker.Location, route *shuttletracker.Route) Prediction {
	// Find the index of the closest point to this shuttle's last known location
	index := ClosestPointTo(lastUpdate.Latitude, lastUpdate.Longitude, route)

	// Figures out whether or not the shuttle is currently at a stop
	atStop := false
	var thisStop = -1
	for i, stop := range route.StopIDs {
		stopIndex := ClosestPointTo(stop.Latitude, stop.Longitude, route)
		if route.Points[index] == route.StopIDs[i]{
			atStop = true
			thisStop = i
			break
		}
	}

	// Find the amount of time that has passed since the last update was received, and given that,
	// the distance the shuttle is predicted to have travelled
	secondsSinceUpdate := time.Since(lastUpdate.Time).Seconds()
	predictedDistance := secondsSinceUpdate * lastUpdate.Speed

	// Iterate over each point in the route in order, summing the distance between each point,
	// and stop when the predicted distance has elapsed
	elapsedDistance := 0.0
	for elapsedDistance < predictedDistance {
		prevIndex := index
		index++
		if index >= len(route.Points) {
			index = 0
		}
		elapsedDistance += DistanceBetween(route.Points[prevIndex], route.Points[index])
	}
	return Prediction{VehicleID: vehicle.ID, Point: route.Points[index], Index: index}
}
