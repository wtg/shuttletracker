package smooth

import (
	"math"
	"time"
	//"fmt"
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

//Returns the index ms.Stops() of the closest stop up ahead on the route (from ms.Stops()) to the given latitude and longitude coordinates
func ClosestApproachingStop(currentIndex int, route *shuttletracker.Route, ms shuttletracker.ModelService) int{
	minDistance := math.MaxFloat64
	minDistanceStopIndex := 0
	i := 0
	stops, err := ms.Stops()
		if err != nil {
			log.WithError(err).Errorf("Unable to retrieve stops")
			return -1
		}
	for i < len(route.StopIDs) {
		j := 0 
		
		for j < len(stops) {
			if (route.StopIDs[i] == stops[j].ID) {
				stopindex := ClosestPointTo(stops[j].Latitude, stops[j].Longitude, route)
				dist := DistanceToPointAlongRoute(currentIndex, stopindex , route)
				//dist := DistanceBetween(shuttletracker.Point{Latitude: stops[j].Latitude, Longitude: stops[j].Longitude} ,route.Points[currentIndex] ) 
				if dist < minDistance {
					minDistance = dist
					minDistanceStopIndex = j
					break
				}
			}
			j += 1
		}
		i += 1
	}
	log.Info("Min distance Approaching      " , minDistance)
	return  minDistanceStopIndex
}

//Returns the index in ms.Stops() of the closest stop up ahead on the route (from ms.Stops()) to the given latitude and longitude coordinates
func ClosestDepartingStop(currentIndex int, route *shuttletracker.Route, ms shuttletracker.ModelService) int{
	minDistance := math.MaxFloat64
	minDistanceStopIndex := 0
	i := 0
	stops, err := ms.Stops()
		if err != nil {
			log.WithError(err).Errorf("Unable to retrieve stops")
			return -1
		}
	for i < len(route.StopIDs) {
		j := 0 
		
		for j < len(stops) {
			if (route.StopIDs[i] == stops[j].ID) {
				stopindex := ClosestPointTo(stops[j].Latitude, stops[j].Longitude, route)
				dist := DistanceToPointAlongRoute( stopindex , currentIndex,  route)
				//dist := DistanceBetween(shuttletracker.Point{Latitude: stops[j].Latitude, Longitude: stops[j].Longitude} ,route.Points[currentIndex] ) 
				if dist < minDistance {
					minDistance = dist
					minDistanceStopIndex = j
					break
				}
			}
			j += 1
		}
		i += 1
	}
	log.Info("Min distance Departing      "  , minDistance)
	return  minDistanceStopIndex
}


func DistanceToPointAlongRoute(currentindex int, stopindex int, route *shuttletracker.Route) float64 {
	elapsedDistance := 0.0
	index := currentindex
	for index != stopindex{
		prevIndex := index
		index++
		if index >= len(route.Points) {
			index = 0
		}
		elapsedDistance += DistanceBetween(route.Points[prevIndex], route.Points[index])
	}
	return elapsedDistance
}
func acceleration(velocity, finalVelocity float64, distance float64) float64 {
	return (math.Pow(finalVelocity, 2) - math.Pow(velocity, 2))  / (2.0 * distance)
}


// Naive algorithm to predict the position a shuttle is at, given the last update received
// Returns the index of the point the shuttle would be at on its route
// TODO: More factors this algorithm should consider: shuttle's proximity to a stop, whether
// the shuttle is going around a sharp turn, etc.
func NaivePredictPosition(vehicle *shuttletracker.Vehicle, lastUpdate *shuttletracker.Location, route *shuttletracker.Route, ms shuttletracker.ModelService) Prediction {
	// Find the index of the closest point to this shuttle's last known location
	log.Info("In NAIVEPREDICTPOSITION \n")
	stops,_ := ms.Stops()
	index := ClosestPointTo(lastUpdate.Latitude, lastUpdate.Longitude, route)

	approachingStopIndex := ClosestApproachingStop(index, route,  ms)
	departingStopIndex :=  ClosestDepartingStop(index, route, ms)
	approachingStopIndex  = ClosestPointTo( stops[approachingStopIndex].Latitude, stops[approachingStopIndex].Longitude, route)
	departingStopIndex  =  ClosestPointTo(stops[departingStopIndex].Latitude, stops[departingStopIndex].Longitude, route)


	//log.Info("approaching stop     ", approachingStopIndex)
	//log.Info("Departing stop       ", departingStopIndex)

	approachingDistance  :=  DistanceToPointAlongRoute(index, approachingStopIndex, route) //meters
	departingDistance :=  DistanceToPointAlongRoute(departingStopIndex, index, route) //meters
	
	log.Info("approaching stop distance     ", approachingDistance)
	log.Info("Departing stop  distance      ", departingDistance)
	
	// Find the amount of time that has passed since the last update was received, and given that,
	// the distance the shuttle is predicted to have travelled
	secondsSinceUpdate := time.Since(lastUpdate.Time).Seconds()
	predictedDistance := secondsSinceUpdate * lastUpdate.Speed

/*
	//if close and approaching a stop:
	if headingToward && distanceToStop < 15 {
		a := -1 * acceleration(lastUpdate.Speed, 0,  distanceToStop)
		predictedDistance = predictedDistance + (.5)*a*math.Pow(secondsSinceUpdate, 2)
	}
	
	// Iterate over each point in the route in order, summing the distance between each point,
	// and stop when the predicted distance has elapsed
	*/
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
