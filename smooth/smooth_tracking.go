package smooth

import (
	"math"
	"time"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
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

	// Iterate over each point in the route in order, summing the distance between each point,
	// and stop when the predicted distance has elapsed
	elapsedDistance := 0.0
	angle := 0.0
	prevAngle := 0.0
	prevDistance := 0.0
	for elapsedDistance < predictedDistance {
		prevAngle = angle
		prevIndex := index
		prevDistance = elapsedDistance
		index++
		if index >= len(route.Points) {
			index = 0
		}
		elapsedDistance += DistanceBetween(route.Points[prevIndex], route.Points[index])
		angle = AngleBetween(route.Points[prevIndex], route.Points[index])

		changeInAngle := math.Abs(math.Mod(angle, 360.0) - math.Mod(prevAngle, 360.0))
		changeInDistance := elapsedDistance - prevDistance

		if changeInAngle > 50 && changeInAngle < 100 && changeInDistance > 1 { // sharp turn and distance traveled
			log.Debugf("Predicting a sharp turn...")
			log.Debugf("Change in angle is %f for vehicle %d", changeInAngle, vehicle.ID)
			log.Debugf("Current Speed is %f for vehicle %d", lastUpdate.Speed, vehicle.ID)
			log.Debugf("Change in distance is %f", changeInDistance)

			// Plan is to drop speed to 6-8 mph when a sharp turn is detected
			// Avg Error is 600-700 meters without this change
			// Running the webapp 3 times found the shuttle speeds to predominantly be
			// 25, 20, 15, 10-12, or 6-8 MPH
			// Assuming 6-8 mph is the turning zone (hard to verify with updates being spotty)

			// Concern
			// Thing is what if a prediction is made and then it doesn't reflect
			// because the angle doesn't isn't large until it actually already starts the turn and slowed down...

			// Also if the last update's speed is used to calculate the distance to be travelled, it
			// means that we can change the distance we need to accumulate before the loop ends
			// but not directly change the speed...
			// I guess we'll add a certain amount of elapsedDistance

			// Say lastUpdate.Speed = 15 mph, and assuming turn is detected and lasts for 2 seconds
			// at 8 MPH, we add 7mph * 2 secs = 6.26 meters, so we add 6.26 meters to elapsedDistance
			// Hard to measure it's impacts unless we run for the entire route as the avg error is
			// 600-700 meters... The error is so large because updates are so rare...
			elapsedDistance += (lastUpdate.Speed - 3.575) * 2 // 8 mph in meters and 2 for # of seconds
		}
	}

	return Prediction{VehicleID: vehicle.ID, Point: route.Points[index], Index: index, Angle: angle}
}
