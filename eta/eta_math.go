package eta

import (
	//"errors"
	"math"
	//"sync"
	//"time"

	// "github.com/wcharczuk/go-chart"

	"github.com/wtg/shuttletracker"
	//"github.com/wtg/shuttletracker/log"
	//"github.com/wtg/shuttletracker/updater"
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

func calculateRouteDistance(route *shuttletracker.Route) float64 {
	totalDistance := 0.0
	for i, p1 := range route.Points {
		if i == len(route.Points)-1 {
			break
		}
		p2 := route.Points[i+1]
		totalDistance += distanceBetween(p1, p2)
	}
	return totalDistance
}

type locationDistance struct {
	loc   *shuttletracker.Location
	dist  float64
	index int
}

func calculateDistance(lds []locationDistance) float64 {
	total := 0.0
	if len(lds) < 2 {
		return total
	}
	for i := range lds[1:] {
		l1 := lds[i].loc
		p1 := shuttletracker.Point{Latitude: l1.Latitude, Longitude: l1.Longitude}
		l2 := lds[i+1].loc
		p2 := shuttletracker.Point{Latitude: l2.Latitude, Longitude: l2.Longitude}
		total += distanceBetween(p1, p2)
	}
	return total
}

// Loop over all points in route and return the route's points that are closest
// to the provided point. We can't just find the distances to all points and then
// take the two points with smallest distances, since they might not be next to
// each other on the route (and therefore are not representative of the route).
func findClosestLine(point shuttletracker.Point, route *shuttletracker.Route) (p1, p2 shuttletracker.Point) {
	if len(route.Points) < 2 {
		return
	}

	// this is the sum of distances from input point to two consecutive points on route.
	// we want to minimize it.
	totalDistance := math.Inf(1)

	for i := range route.Points[1:] {
		tempP1 := route.Points[i]
		tempP2 := route.Points[i+1]
		d1 := distanceBetween(point, tempP1)
		d2 := distanceBetween(point, tempP2)
		d := d1 + d2
		if d < totalDistance {
			p1 = tempP1
			p2 = tempP2
			totalDistance = d
		}
	}

	return
}

// for each location, return the index of the provided point which it is closest to.
// this can be used e.g. to figure out the order in which a track traverses a list of stops.
func findMinimumDistanceIndices(points []shuttletracker.Point, locs []*shuttletracker.Location) []int {
	indices := make([]int, len(locs))
	for i, loc := range locs {
		locPoint := shuttletracker.Point{Latitude: loc.Latitude, Longitude: loc.Longitude}
		minDistance := math.Inf(1)
		var minIndex int
		for j, p := range points {
			d := distanceBetween(p, locPoint)
			if d < minDistance {
				minIndex = j
				minDistance = d
			}
		}
		indices[i] = minIndex
	}
	return indices
}

// lat/lon. returns bearing in degrees
func findInitialBearing(p1, p2 shuttletracker.Point) float64 {
	lat1Rad := toRadians(p1.Latitude)
	lon1Rad := toRadians(p1.Longitude)
	lat2Rad := toRadians(p2.Latitude)
	lon2Rad := toRadians(p2.Longitude)
	lonDiff := lon2Rad - lon1Rad

	y := math.Sin(lonDiff) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(lonDiff)
	return math.Atan2(y, x) / (math.Pi / 180)
}

// cross-track distance. see http://www.movable-type.co.uk/scripts/latlong.html
// Sign of returned value indicates which side of route the point is on.
func crossTrackDistance(p shuttletracker.Point, route *shuttletracker.Route) float64 {
	// first find two points that define the line
	p1, p2 := findClosestLine(p, route)

	// find angular distance from first line point to input point
	angDist := distanceBetween(p1, p) / earthRadius

	// find bearing from first line point to input point
	b1 := findInitialBearing(p1, p)

	// find bearing from first line point to second line point
	b2 := findInitialBearing(p1, p2)

	return math.Asin(math.Sin(angDist)*math.Sin(b1-b2)) * earthRadius
}
