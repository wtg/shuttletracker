package smooth

import (
	"math"
	"time"

	"github.com/wtg/shuttletracker"
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

	return brng - 45
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
	// Dictionary of hard-coded light stops and intersections

	numLights := 22
	m := make(map[int64][]float64)
	m[1] = []float64{42.730775, -73.677216}
	m[2] = []float64{42.732904, -73.682553}
	m[3] = []float64{42.733140, -73.683505}
	m[4] = []float64{42.733779, -73.683293}
	m[5] = []float64{42.733708, -73.685683}
	m[6] = []float64{42.732102, -73.686238}
	m[7] = []float64{42.729597, -73.686808}
	m[8] = []float64{42.728325, -73.687088}
	m[9] = []float64{42.727317, -73.687340}
	m[10] = []float64{42.722734, -73.679778}
	m[11] = []float64{42.726857, -73.678091}
	m[12] = []float64{42.733992, -73.682248}
	m[13] = []float64{42.732400, -73.671711}
	m[14] = []float64{42.736453, -73.670587}
	m[15] = []float64{42.737615, -73.670265}
	m[16] = []float64{42.738269, -73.670081}
	m[17] = []float64{42.737730, -73.666570}
	m[18] = []float64{42.735967, -73.667075}
	m[19] = []float64{42.730862, -73.667257}
	m[20] = []float64{42.726621, -73.666819}
	m[21] = []float64{42.724851, -73.673798}
	m[22] = []float64{42.729077, -73.672603}

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
		angle = AngleBetween(route.Points[prevIndex], route.Points[index]) + 45

		changeInAngle := math.Abs(math.Mod(angle, 360.0) - math.Mod(prevAngle, 360.0))
		changeInDistance := elapsedDistance - prevDistance

		if changeInAngle > 50 && changeInAngle < 100 && changeInDistance > 1 && lastUpdate.Speed > 3.6 {
			elapsedDistance += (lastUpdate.Speed - 3.575) * 2 // 8 mph in meters and 2 for # of seconds
		}
	}

	return Prediction{VehicleID: vehicle.ID, Point: route.Points[index], Index: index, Angle: angle}
}
