package eta

import (
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"strings"
	"math"
	"os"
	"fmt"
	"sync"

	"github.com/wcharczuk/go-chart"
	"github.com/spf13/viper"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/updater"
)

const osrmBaseURL = "http://127.0.0.1:8080/api/osrm"
const earthRadius = 6371000.0 // meters

type ETAManager struct {
	ms shuttletracker.ModelService
	locChan chan *shuttletracker.Location

	sm *sync.Mutex
	subscribers []func(VehicleETA)

	em *sync.Mutex
	etas []VehicleETA
}

type VehicleETA struct {
	VehicleID int64
	RouteID int64
	StopETAs map[int64]time.Time
}

type osrmRouteResp struct {
	Routes []struct {
		Duration float64
	}
}

type point struct {
	lat float64
	lng float64
}

type Config struct {
	DataFeed       string
	UpdateInterval string
}

// NewManager creates an ETAManager subscribed to Location updates from Updater.
func NewManager(cfg Config, ms shuttletracker.ModelService, updater *updater.Updater) (*ETAManager, error) {
	em := &ETAManager{
		ms: ms,
		locChan: make(chan *shuttletracker.Location),
		sm: &sync.Mutex{},
		subscribers: []func(VehicleETA){},
		em: &sync.Mutex{},
		etas: []VehicleETA{},
	}

	// subscribe to new Locations with Updater
	updater.Subscribe(em.locationSubscriber)

	return em, nil
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		// UpdateInterval: "10s",
		// DataFeed:       "https://shuttles.rpi.edu/datafeed",
	}
	// v.SetDefault("updater.updateinterval", cfg.UpdateInterval)
	// v.SetDefault("updater.datafeed", cfg.DataFeed)
	return cfg
}

func (em *ETAManager) locationSubscriber(loc *shuttletracker.Location) {
	em.locChan <- loc
}

func (em *ETAManager) Run() {
	for {
		select {
		case loc := <-em.locChan:
			em.handleNewLocation(loc)
		}
	}
}

func (em *ETAManager) handleNewLocation(loc *shuttletracker.Location) {
	log.Infof("ETAManager got location: %+v", loc)
}

func (em *ETAManager) Subscribe(sub func(VehicleETA)) {
	em.sm.Lock()
	em.subscribers = append(em.subscribers, sub)
	em.sm.Unlock()
}

func (em *ETAManager) ETAs() []VehicleETA {
	em.em.Lock()
	etas := make([]VehicleETA, len(em.etas))
	copy(etas, em.etas)
	em.em.Unlock()
	return etas
}

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
		haversine(lat2Rad - lat1Rad) +
		math.Cos(lat1Rad) * math.Cos(lat2Rad) *
		haversine(lon2Rad - lon1Rad)))
}

func closestPointIndexOnRoute(stop *shuttletracker.Stop, route []shuttletracker.Point) (int, error) {
	minDistance := math.Inf(0)
	var minIndex int
	for i, p1 := range route {
		var p2 shuttletracker.Point
		if i < len(route) - 1 {
			p2 = route[i+1]
		} else {
			p2 = route[0]
		}

		// find distance from stop to line defined by two points
		d := math.Abs((p2.Latitude - p1.Latitude) * stop.Longitude -
			(p2.Longitude - p1.Longitude) * stop.Latitude +
			(p2.Longitude * p1.Latitude) - (p2.Latitude * p1.Longitude)) /
			math.Sqrt(math.Pow(p2.Latitude - p1.Latitude, 2) + math.Pow(p2.Longitude - p1.Longitude, 2))

		// d := distanceBetween(p, stopPoint)
		if d < minDistance {
			minDistance = d
			minIndex = i
		}
	}

	return minIndex, nil
}

func distanceBetweenStopsAlongRoute(stop1, stop2 *shuttletracker.Stop, route *shuttletracker.Route) float64 {
	// p1, p2, err := closestPointsOnRoute(stop1, route)
	// if err != nil {
	// 	log.WithError(err).Error("unable to find closest point on route")
	// 	return -1
	// }
	// _ = p1
	// _ = p2
	return 0
}

// Loop over all points in all segments and return the index of the segment containing the closest point.
func snapPointToSegmentIndex(point shuttletracker.Point, segments [][]shuttletracker.Point) int {
	minDistance := math.Inf(0)
	var minIndex int

	for i, segment := range segments {
		// log.Info(segment)
		for _, p := range segment {
			d := distanceBetween(point, p)
			if d < minDistance {
				minIndex = i
			}
		}
	}

	return minIndex
}

func (em *ETAManager) snapLocationsToRoads(locations []*shuttletracker.Location) ([]point, error) {
	coords := strings.Builder{}
	for i, loc := range locations {
		coords.WriteString(strconv.FormatFloat(loc.Longitude, 'f', -1, 64))
		coords.WriteByte(',')
		coords.WriteString(strconv.FormatFloat(loc.Latitude, 'f', -1, 64))
		if i != len(locations) - 1 {
			coords.WriteByte(';')
		}
	}

	c := &http.Client{Timeout: 5 * time.Second}
	resp, err := c.Get(osrmBaseURL + "/match/v1/car/" + coords.String())
	if err != nil {
		return nil, err
	}

	osrmMatch := osrmRouteResp{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&osrmMatch)
	if err != nil {
		return nil, err
	}
	log.Infof("%+v", osrmMatch)
	return nil, nil
}

func calculateRouteDistance(route *shuttletracker.Route) float64 {
	totalDistance := 0.0
	for i, p1 := range route.Points {
		if i == len(route.Points) - 1 {
			break
		}
		p2 := route.Points[i+1]
		totalDistance += distanceBetween(p1, p2)
	}
	return totalDistance
}

type locationDistance struct {
	loc *shuttletracker.Location
	dist float64
	index int
}

func findNextMinimum(locDists []locationDistance) *locationDistance {
	for _, ld := range locDists {
		if ld.dist < 0.03 {
			return &ld
		}
	}
	return nil
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

// determines if a track has traversed a route in the order of its stops. if a stop was missing,
// it is ignored. this is because of stops that are close together often not having enough locations.
// only two stops may be dropped before we bail out.
func (em *ETAManager) stopsVisitedInOrder(route *shuttletracker.Route, locDists []locationDistance) (bool, error) {
	// return true
	stopPoints := make([]shuttletracker.Point, len(route.StopIDs))
	for i, stopID := range route.StopIDs {
		stop, err := em.ms.Stop(stopID)
		if err != nil {
			return false, err
		}
		p := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}
		stopPoints[i] = p
	}

	locs := make([]*shuttletracker.Location, len(locDists))
	for i, ld := range locDists {
		locs[i] = ld.loc
	}

	// associate locations with nearest stops and then turn indices into stop IDs
	minIndices := findMinimumDistanceIndices(stopPoints, locs)
	stopIDs := make([]int64, len(minIndices))
	for i, min := range minIndices {
		stopIDs[i] = route.StopIDs[min]
	}
	// log.Debugf("%+v", stopIDs)

	// if a stop is missing from the track, ignore it. this can happen if stops are close
	// together because the location data is sometimes infrequent.
	desiredStopIDs := []int64{}
	for _, desiredID := range route.StopIDs {
		for _, stopID := range stopIDs {
			if desiredID == stopID {
				desiredStopIDs = append(desiredStopIDs, desiredID)
				break
			}
		}
	}
	// only allow dropping at most two stops
	if len(route.StopIDs) - len(desiredStopIDs) > 2 {
		return false, nil
	}

	// were the stops visited in order?
	for i := 0; i < len(desiredStopIDs); i++ {
		found := false
		for j := i; j < len(stopIDs); j++ {
			if desiredStopIDs[i] == stopIDs[j] {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}


	return true, nil
}

// lat/lon. returns bearing in degrees
func findInitialBearing(p1, p2 shuttletracker.Point) float64 {
	lat1Rad := toRadians(p1.Latitude)
	lon1Rad := toRadians(p1.Longitude)
	lat2Rad := toRadians(p2.Latitude)
	lon2Rad := toRadians(p2.Longitude)
	lonDiff := lon2Rad - lon1Rad

	y := math.Sin(lonDiff) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad) * math.Sin(lat2Rad) - math.Sin(lat1Rad) * math.Cos(lat2Rad) * math.Cos(lonDiff)
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

	return math.Asin(math.Sin(angDist) * math.Sin(b1 - b2)) * earthRadius
}

func (em *ETAManager) findTracks(locDists []locationDistance, route *shuttletracker.Route) ([][]locationDistance, error) {
	tracks := [][]locationDistance{}

	for i := 0; i < len(locDists); i++ {
		// log.Debugf("tracks i %d", i)
		// go until we find a departure from the first stop (distance is small)
		ld := locDists[i]
		if (ld.dist > 30) {
			continue
		}

		// found it. now follow the track around the route until we get to the first stop again.
		track := []locationDistance{ld}
		i++
		for ; i < len(locDists)-1; i++ {
			ld = locDists[i]

			// add location to track
			track = append(track, ld)

			// end track if location is back at the initial stop
			if (ld.dist < 30) {
				break
			}
		}

		// check that all locations in track are on the route
		onRoute := true
		for _, ld := range track {
			if ld.loc.RouteID == nil || *ld.loc.RouteID != route.ID {
				onRoute = false
				break
			}
		}
		if !onRoute {
			continue
		}

		// see if this track is valid
		// does it have at least five locations?
		if len(track) < 5 {
			continue
		}
		// at least five minutes elapsed?
		if track[len(track)-1].loc.Time.Sub(track[0].loc.Time) < 5 * time.Minute {
			continue
		}
		// total distance traveled is at least 75% of route length? no more than
		// 110% of route length? (it is unlikely that a track is longer than a route
		// since a route has many more points than a track—a track essentially cuts corners.)
		routeLength := calculateRouteDistance(route)
		d := calculateDistance(track)
		if d < routeLength * 0.75 || d > routeLength * 1.1 {
			continue
		}

		// no point on the track more than 100 m from the route?
		nearRoute := true
		for _, ld := range track {
			p := shuttletracker.Point{Latitude: ld.loc.Latitude, Longitude: ld.loc.Longitude}
			if math.Abs(crossTrackDistance(p, route)) > 100 {
				nearRoute = false
				break
			}
		}
		if !nearRoute {
			continue
		}

		// stops visited in correct order?
		inOrder, err := em.stopsVisitedInOrder(route, track)
		if err != nil {
			return nil, err
		}
		if !inOrder {
			continue
		}

		// this track looks good
		tracks = append(tracks, track)
	}
	return tracks, nil
}

func (em *ETAManager) findRouteLoops(route *shuttletracker.Route) ([][]*shuttletracker.Location, error) {
	vehicles, err := em.ms.Vehicles()
	if err != nil {
		return nil, err
	}

	// this is the first stop on the route. we consider a loop to be a departure from
	// this stop followed by an eventual arrival at this stop.
	if len(route.StopIDs) == 0 {
		return nil, nil
	}
	stopID := route.StopIDs[0]
	stop, err := em.ms.Stop(stopID)
	if err != nil {
		return nil, err
	}
	stopPoint := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}

	graph := chart.Chart{Series: []chart.Series{}}
	xVals := []float64{}
	yVals := []float64{}
	for _, point := range route.Points {
		xVals = append(xVals, point.Longitude)
		yVals = append(yVals, point.Latitude)
	}
	graph.Series = append(graph.Series, chart.ContinuousSeries{
		Name: fmt.Sprintf("%s", route.Name),
		XValues: xVals,
		YValues: yVals,
	})

	tracks := [][]locationDistance{}
	for _, vehicle := range vehicles {
		locations, err := em.ms.LocationsSince(vehicle.ID, time.Now().Add(-time.Hour*24*30))
		if err != nil {
			return nil, err
		}

		// reverse locations so oldest is first
		for left, right := 0, len(locations)-1; left < right; left, right = left+1, right-1 {
			locations[left], locations[right] = locations[right], locations[left]
		}

		// associate each location with distance to first stop on route
		locDistances := []locationDistance{}
		for i, loc1 := range locations {
			locPoint := shuttletracker.Point{Latitude: loc1.Latitude, Longitude: loc1.Longitude}
			d := distanceBetween(stopPoint, locPoint)
			locDistances = append(locDistances, locationDistance{loc: loc1, dist: d, index: i})
		}

		vehicleTracks, err := em.findTracks(locDistances, route)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, vehicleTracks...)
	}

	totalTime := 0.0
	for _, track := range tracks {
		xVals := []float64{}
		yVals := []float64{}
		for _, loc := range track {
			xVals = append(xVals, loc.loc.Longitude)
			yVals = append(yVals, loc.loc.Latitude)
		}
		graph.Series = append(graph.Series, chart.ContinuousSeries{
			// Name: fmt.Sprintf("%d", i),
			XValues: xVals,
			YValues: yVals,
		})

		totalTime += track[len(track)-1].loc.Time.Sub(track[0].loc.Time).Seconds()
	}
	avgTime := totalTime / float64(len(tracks))
	log.Debugf("avg: %f min", avgTime/60.0)

	graph.Elements = []chart.Renderable{chart.Legend(&graph)}
	f, _ := os.Create(fmt.Sprintf("plot-%s.png", route.Name))
	graph.Render(chart.PNG, f)
	f.Close()

	return nil, nil
}

func (em *ETAManager) calculateInitialETAs() {
	routes, err := em.ms.Routes()
	if err != nil {
		log.WithError(err).Error("unable to get routes")
		return
	}
	for _, route := range routes {
		routeDistance := calculateRouteDistance(route)
		log.Debugf("%s %f", route.Name, routeDistance)
		em.findRouteLoops(route)
	}
}

func (em *ETAManager) calculateETAs(vehicle *shuttletracker.Vehicle) {
	log.Infof("calculating ETA for vehicle %+v", vehicle)
	location, err := em.ms.LatestLocation(vehicle.ID)
	if err != nil {
		log.WithError(err).Error("unable to get latest location")
		return
	}

	lonStr := strconv.FormatFloat(location.Longitude, 'f', -1, 64)
	latStr := strconv.FormatFloat(location.Latitude, 'f', -1, 64)

	c := &http.Client{Timeout: 5 * time.Second}

	stops, err := em.ms.Stops()
	for _, stop := range stops {
		stopLonStr := strconv.FormatFloat(stop.Longitude, 'f', -1, 64)
		stopLatStr := strconv.FormatFloat(stop.Latitude, 'f', -1, 64)

		coords := lonStr + "," + latStr + ";" + stopLonStr + "," + stopLatStr

		resp, err := c.Get(osrmBaseURL + "/route/v1/car/" + coords)
		if err != nil {
			log.WithError(err).Error("unable to get OSRM route")
			continue
		}

		osrmRoute := osrmRouteResp{}
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&osrmRoute)
		if err != nil {
			log.WithError(err).Error("unable to decode")
			continue
		}
		log.Infof("%+v", osrmRoute)
	}
}
