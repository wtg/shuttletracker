package smooth

import (
	"math"
	"time"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/eta"
)


type Prediction struct {
	VehicleID int64
	Point     shuttletracker.Point
	Index     int
	Angle     float64
}

const earthRadius = 6371000.0 // meters

func haversine(theta float64) float64 {
	return (1 - math.Cos(theta)) / 2
}

func toRadians(n float64) float64 {
	return n * math.Pi / 180
}

func toDegrees(angle float64) float64 {
	return angle * 180 / math.Pi
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

func AngleBetween(p1, p2 shuttletracker.Point) float64 {
	radLat1 := toRadians(p1.Latitude)
	radLng1 := toRadians(p1.Longitude)
	radLat2 := toRadians(p2.Latitude)
	radLng2 := toRadians(p2.Longitude)

	deltaLongitude := (radLng2 - radLng1)
	y := math.Sin(deltaLongitude) * math.Cos(radLat2)
	x := math.Cos(radLat1)*math.Sin(radLat2) - math.Sin(radLat1)*math.Cos(radLat2)*math.Cos(deltaLongitude)

	brng := math.Atan2(y, x)
	brng = toDegrees(brng)
	brng = math.Mod(brng+360, 360)
	brng = 360 - brng // Convert to counter-clockwise

	return brng
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
// TODO: More factors this algorithm should consider: shuttle's proximity to a stop, whether
// the shuttle is going around a sharp turn, etc.
func NaivePredictPosition(vehicle *shuttletracker.Vehicle, lastUpdate *shuttletracker.Location, route *shuttletracker.Route) Prediction {
	// Find the index of the closest point to this shuttle's last known location
	index := ClosestPointTo(lastUpdate.Latitude, lastUpdate.Longitude, route)

	// Find the amount of time that has passed since the last update was received, and given that,
	// the distance the shuttle is predicted to have travelled
	secondsSinceUpdate := time.Since(lastUpdate.Time).Seconds()
	predictedDistance := secondsSinceUpdate * lastUpdate.Speed
	
	// if we are within the radius of a stop.
	/*
	em = NewManager(ms, updater) (*ETAManager, error)
	lastDepartureTrack, err1 := em.GetLastDepartureTrack(vehicle, route)

	locIndices, err2 := em.DetermineNearestStopIndicesAlongRoute(lastDepartureTrack, route)
	
	if err1 == nil || err2 == nil { 

	}
	else {
		closeToStop = false
	}
	*/

	// Iterate over each point in the route in order, summing the distance between each point,
	// and stop when the predicted distance has elapsed
	elapsedDistance := 0.0
	angle := 0.0
	for elapsedDistance < predictedDistance {
		prevIndex := index
		index++
		if index >= len(route.Points) {
			index = 0
		}
		elapsedDistance += DistanceBetween(route.Points[prevIndex], route.Points[index])
		angle = AngleBetween(route.Points[prevIndex], route.Points[index])
	}
	return Prediction{VehicleID: vehicle.ID, Point: route.Points[index], Index: index, Angle: angle}
}
