package smooth

import (
	"math"
	"time"

	"github.com/wtg/shuttletracker"
)

const earthRadius = 6371000.0 // meters

func haversine(theta float64) float64 {
	return (1 - math.Cos(theta)) / 2
}

func toRadians(n float64) float64 {
	return n * math.Pi / 180
}

func distanceBetween(p1, p2 shuttletracker.Point) float64 {
	lat1Rad := toRadians(p1.Latitude)
	lon1Rad := toRadians(p1.Longitude)
	lat2Rad := toRadians(p2.Latitude)
	lon2Rad := toRadians(p2.Longitude)

	return 2 * earthRadius * math.Asin(math.Sqrt(
		haversine(lat2Rad-lat1Rad)+
			math.Cos(lat1Rad)*math.Cos(lat2Rad)*
				haversine(lon2Rad-lon1Rad)))
}

func estimatePosition(vehicle *shuttletracker.Vehicle, lastUpdate *shuttletracker.Location, route *shuttletracker.Route) int {
	index := 0
	minDistance := 0.0
	for i, point := range route.Points {
		distance := distanceBetween(point, shuttletracker.Point{Latitude: lastUpdate.Latitude, Longitude: lastUpdate.Longitude})
		if distance < minDistance {
			minDistance = distance
			index = i
		}
	}
	secondsSinceUpdate := time.Since(lastUpdate.Time).Seconds()
	predictedDistance := secondsSinceUpdate * lastUpdate.Speed
	elapsedDistance := 0.0
	for elapsedDistance < predictedDistance {
		prevIndex := index
		index++
		if index >= len(route.Points) {
			index = 0
		}
		elapsedDistance += distanceBetween(route.Points[prevIndex], route.Points[index])
	}
	return index
}
