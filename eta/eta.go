package eta

import (
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"strings"
	"math"
	// "os"
	// "fmt"
	"sync"
	"errors"

	// "github.com/wcharczuk/go-chart"
	"github.com/spf13/viper"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/updater"
)

const osrmBaseURL = "http://127.0.0.1:8080/api/osrm"
const earthRadius = 6371000.0 // meters

type ETAManager struct {
	ms shuttletracker.ModelService
	etaChan chan *VehicleETA
	etas map[int64]*VehicleETA
	etasReqChan chan chan map[int64]*VehicleETA

	sm *sync.Mutex
	subscribers []func(VehicleETA)
}

type VehicleETA struct {
	VehicleID int64 `json:"vehicle_id"`
	RouteID int64 `json:"route_id"`
	StopETAs []StopETA `json:"stop_etas"`
	Updated time.Time `json:"updated"`
}

type StopETA struct {
	StopID int64 `json:"stop_id"`
	ETA time.Time `json:"eta"`
	Arriving bool `json:"arriving"`
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
		etaChan: make(chan *VehicleETA),
		etas: map[int64]*VehicleETA{},
		etasReqChan: make(chan chan map[int64]*VehicleETA),
		sm: &sync.Mutex{},
		subscribers: []func(VehicleETA){},
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

// This gets new Locations from Updater. As soon as this happens, we'll
// determine new ETAs for the vehicle in another goroutine.
func (em *ETAManager) locationSubscriber(loc *shuttletracker.Location) {
	go em.handleNewLocation(loc)
}

// Run is in charge of managing all of the state inside of ETAManager.
func (em *ETAManager) Run() {
	ticker := time.Tick(time.Minute)
	for {
		select {
		case eta := <-em.etaChan:
			em.handleNewETA(eta)
		case etasReplyChan := <-em.etasReqChan:
			em.processETAsRequest(etasReplyChan)
		case <-ticker:
			em.cleanup()
		}
	}
}

func (em *ETAManager) handleNewLocation(loc *shuttletracker.Location) {
	if loc.VehicleID == nil {
		// can't do anything...
		return
	}
	vehicleID := *loc.VehicleID
	em.calculateVehicleETAs(vehicleID)
}

func (em *ETAManager) handleNewETA(eta *VehicleETA) {
	em.etas[eta.VehicleID] = eta

	// notify subscribers
	em.sm.Lock()
	for _, sub := range em.subscribers {
		sub(*eta)
	}
	em.sm.Unlock()
}

// spit out all current ETAs over the provided channel
func (em *ETAManager) processETAsRequest(c chan map[int64]*VehicleETA) {
	etas := map[int64]*VehicleETA{}
	for k, v := range em.etas {
		etas[k] = v
	}

	c <- etas
}

// iterate over all ETAs and remove those that have expired.
func (em *ETAManager) cleanup() {
	log.Debug("ETAManager cleanup")
	now := time.Now()
	for vehicleID, vehicleETA := range em.etas {
		stopETAs := vehicleETA.StopETAs
		for i := len(stopETAs) - 1; i >= 0; i-- {
			stopETA := stopETAs[i]
			if now.After(stopETA.ETA) {
				stopETAs = append(stopETAs[:i], stopETAs[i+1:]...)
			}
		}
		if len(stopETAs) == 0 {
			delete(em.etas, vehicleID)
		} else {
			vehicleETA.StopETAs = stopETAs
		}
	}
}

func (em *ETAManager) Subscribe(sub func(VehicleETA)) {
	em.sm.Lock()
	em.subscribers = append(em.subscribers, sub)
	em.sm.Unlock()
}

// this can be called by anyone to get ETAManager's current view of vehcile ETAs
func (em *ETAManager) CurrentETAs() map[int64]*VehicleETA {
	etasChan := make(chan map[int64]*VehicleETA)
	em.etasReqChan <- etasChan
	return <-etasChan
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

func calculatePointsDistance(points []shuttletracker.Point) float64 {
	totalDistance := 0.0
	for i, p1 := range points {
		if i == len(points) - 1 {
			break
		}
		p2 := points[i+1]
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

func (em *ETAManager) determineNearestStopIndicesAlongRoute(track []*shuttletracker.Location, route *shuttletracker.Route) ([]int, error) {
	if len(track) < 2 {
		return nil, errors.New("not enough locations in track")
	}

	// get stops
	stops := make([]*shuttletracker.Stop, len(route.StopIDs))
	for i, stopID := range route.StopIDs {
		stop, err := em.ms.Stop(stopID)
		if err != nil {
			return nil, err
		}
		stops[i] = stop
	}

	stopPoints := make([]shuttletracker.Point, len(route.StopIDs))
	for i, stop := range stops {
		stopPoint := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}
		stopPoints[i] = stopPoint
	}

	// associate locations with nearest stop zones
	stopZones, err := em.findStopZoneIndices(route, track)
	if err != nil {
		return nil, err
	}

	return stopZones, nil
}

func calculateDistanceAlongRoute(locs []*shuttletracker.Location, route *shuttletracker.Route) float64 {
	total := 0.0
	if len(locs) < 2 {
		return total
	}

	// for every subsequent vehicle location, find the distance between the closest route points
	lastRoutePointIdx := 0
	for i := range locs {
		// for every location, find the route points on either side of it
		locPoint := shuttletracker.Point{Latitude: locs[i].Latitude, Longitude: locs[i].Longitude}

		// this is the sum of distances from location point to two consecutive points on route.
		// we want to minimize it.
		totalDistance := math.Inf(1)
		var idx1, idx2 int

		for j := lastRoutePointIdx; j < len(route.Points) - 1; j++ {
			p1 := route.Points[j]
			p2 := route.Points[j+1]

			d1 := distanceBetween(locPoint, p1)
			d2 := distanceBetween(locPoint, p2)
			d := d1 + d2
			
			if d < totalDistance {
				totalDistance = d
				idx1 = j
				idx2 = j+1
			}
		}
		lastRoutePointIdx = idx2
		segmentDistance := calculatePointsDistance(route.Points[idx1:idx2+1])
		total += segmentDistance
		log.Debugf("\n\ntotal: %f\nremaining points: %d\n", total, len(route.Points) - lastRoutePointIdx)
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

// Loop over all points in route and return the index of the route's point that is closest
// to the provided point. Similar to findClosestLine but only returns one point.
func findClosestPointIndex(point shuttletracker.Point, route *shuttletracker.Route) int {
	minIndex := 0
	if len(route.Points) < 1 {
		return minIndex 
	}

	minDistance := math.Inf(1)

	for i := range route.Points {
		p := route.Points[i]
		d := distanceBetween(point, p)
		if d < minDistance {
			minIndex = i
			minDistance = d
		}
	}

	return minIndex
}

// for each location, return the index of the stop zone that it is closest to on the route.
func (em *ETAManager) findStopZoneIndices(route *shuttletracker.Route, locs []*shuttletracker.Location) ([]int, error) {
	stopPoints := make([]shuttletracker.Point, len(route.StopIDs))
	for i, stopID := range route.StopIDs {
		stop, err := em.ms.Stop(stopID)
		if err != nil {
			return nil, err
		}
		p := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}
		stopPoints[i] = p
	}

	indices := make([]int, len(locs))
	minIndex := 0
	for i, loc := range locs {
		locPoint := shuttletracker.Point{Latitude: loc.Latitude, Longitude: loc.Longitude}
		minDistance := math.Inf(1)
		for j := minIndex; j < len(stopPoints); j++ {
			p := stopPoints[j]
			d := distanceBetween(p, locPoint)
			if d < minDistance {
				minIndex = j
				minDistance = d
			}
		}
		indices[i] = minIndex
	}

	return indices, nil
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
func (em *ETAManager) stopsVisitedInOrder(route *shuttletracker.Route, locs []*shuttletracker.Location) (bool, error) {
	stopPoints := make([]shuttletracker.Point, len(route.StopIDs))
	for i, stopID := range route.StopIDs {
		stop, err := em.ms.Stop(stopID)
		if err != nil {
			return false, err
		}
		p := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}
		stopPoints[i] = p
	}

	// associate locations with nearest stops and then turn indices into stop IDs
	minIndices := findMinimumDistanceIndices(stopPoints, locs)
	stopIDs := make([]int64, len(minIndices))
	for i, min := range minIndices {
		stopIDs[i] = route.StopIDs[min]
	}

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

// no point on the track more than 100 m from the route?
func (em *ETAManager) trackNearRoute(track []*shuttletracker.Location, route *shuttletracker.Route) bool {
	for _, loc := range track {
		p := shuttletracker.Point{Latitude: loc.Latitude, Longitude: loc.Longitude}
		if math.Abs(crossTrackDistance(p, route)) > 100 {
			return false
		}
	}
	return true
}

func (em *ETAManager) findTracks(locDists []locationDistance, route *shuttletracker.Route) ([][]locationDistance, error) {
	tracks := [][]locationDistance{}

	for i := 0; i < len(locDists); i++ {
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
		locTrack := make([]*shuttletracker.Location, len(track))
		for i, ld := range track{
			locTrack[i] = ld.loc
		}
		if !em.trackNearRoute(locTrack, route) {
			continue
		}

		// stops visited in correct order?
		inOrder, err := em.stopsVisitedInOrder(route, locTrack)
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
		return nil, errors.New("route doesn't have initial stop")
	}
	stopID := route.StopIDs[0]
	stop, err := em.ms.Stop(stopID)
	if err != nil {
		return nil, err
	}
	stopPoint := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}

	// graph := chart.Chart{Series: []chart.Series{}}
	// xVals := []float64{}
	// yVals := []float64{}
	// for _, point := range route.Points {
	// 	xVals = append(xVals, point.Longitude)
	// 	yVals = append(yVals, point.Latitude)
	// }
	// graph.Series = append(graph.Series, chart.ContinuousSeries{
	// 	Name: fmt.Sprintf("%s", route.Name),
	// 	XValues: xVals,
	// 	YValues: yVals,
	// })

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
		// xVals := []float64{}
		// yVals := []float64{}
		// for _, loc := range track {
		// 	xVals = append(xVals, loc.loc.Longitude)
		// 	yVals = append(yVals, loc.loc.Latitude)
		// }
		// graph.Series = append(graph.Series, chart.ContinuousSeries{
		// 	// Name: fmt.Sprintf("%d", i),
		// 	XValues: xVals,
		// 	YValues: yVals,
		// })

		totalTime += track[len(track)-1].loc.Time.Sub(track[0].loc.Time).Seconds()
	}
	// avgTime := totalTime / float64(len(tracks))

	// graph.Elements = []chart.Renderable{chart.Legend(&graph)}
	// f, _ := os.Create(fmt.Sprintf("plot-%s.png", route.Name))
	// graph.Render(chart.PNG, f)
	// f.Close()

	locs := make([][]*shuttletracker.Location, len(tracks))
	for i, track := range tracks {
		locTrack := make([]*shuttletracker.Location, len(track))
		for j := range track {
			locTrack[j] = track[j].loc
		}
		locs[i] = locTrack
	}

	return locs, nil
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

// WARNING: this assumes that the initial stop only occurs on a route _once_!
func (em *ETAManager) getLastDepartureTrack(vehicle *shuttletracker.Vehicle, route *shuttletracker.Route) ([]*shuttletracker.Location, error) {
	since := time.Now().Add(time.Minute*-30)
	locs, err := em.ms.LocationsSince(vehicle.ID, since)
	if err != nil {
		return nil, err
	}

	// get initial stop on route
	if len(route.StopIDs) == 0 {
		return nil, errors.New("route doesn't have initial stop")
	}
	stopID := route.StopIDs[0]
	stop, err := em.ms.Stop(stopID)
	if err != nil {
		return nil, err
	}
	stopPoint := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}

	if len(locs) < 2 {
		return nil, errors.New("not enough locations")
	}

	i := 0
	for ; i < len(locs); i++ {
		loc := locs[i]
		locPoint := shuttletracker.Point{Latitude: loc.Latitude, Longitude: loc.Longitude}
		d := distanceBetween(stopPoint, locPoint)
		if d < 30 {
			break
		}
	}

	track := make([]*shuttletracker.Location, 0, i)
	for ; i > 0; i-- {
		track = append(track, locs[i-1])
	}
	return track, nil
}

// a "zone" is an area of a route that is closest to a specific stop
func (em *ETAManager) findDurationsByStopZone(locs []*shuttletracker.Location, route *shuttletracker.Route) ([]time.Duration, error) {
	stopZones, err := em.findStopZoneIndices(route, locs)
	if err != nil {
		return nil, err
	}

	durations := make([]time.Duration, len(route.StopIDs))

	entryLoc := locs[0]
	lastZoneIdx := stopZones[0]
	for i := 0; i < len(locs) - 1; i++ {
		zoneIdx := stopZones[i]
		if zoneIdx != lastZoneIdx {
			exitLoc := locs[i]
			duration := exitLoc.Time.Sub(entryLoc.Time)
			durations[zoneIdx] = duration
			entryLoc = exitLoc
			lastZoneIdx = zoneIdx
		}
	}

	return durations, nil
}

// for each point on route, figure out how long it takes a shuttle to reach it
// from the first stop on the route
func (em *ETAManager) determineAverageTravelTimes(route *shuttletracker.Route) ([]time.Duration, error) {
	tracks, err := em.findRouteLoops(route)
	if err != nil {
		return nil, err
	}

	// determine duration of time spent near each stop, grouped by stop index on route
	trackElapseds := make([][]time.Duration, len(route.StopIDs))
	for _, track := range tracks {
		durations, err := em.findDurationsByStopZone(track, route)
		if err != nil {
			return nil, err
		}

		for j, duration := range durations {
			trackElapseds[j] = append(trackElapseds[j], duration)
		}
	}

	// average the durations for each zone
	durations := make([]time.Duration, len(route.Points))
	for i, elapseds := range trackElapseds {
		total := 0.0
		for _, elapsed := range elapseds {
			total += elapsed.Seconds()
		}
		avg := total / float64(len(elapseds))
		durations[i] = time.Duration(avg) * time.Second
	}

	return durations, nil
}

func (em *ETAManager) calculateVehicleETAs(vehicleID int64) {
	// get vehicle info
	vehicle, err := em.ms.Vehicle(vehicleID)
	if err != nil {
		log.WithError(err).Error("unable to get vehicle")
		return
	}

	loc, err := em.ms.LatestLocation(vehicleID)
	if err != nil {
		log.WithError(err).Error("unable to get latest location")
		return
	}

	locPoint := shuttletracker.Point{Latitude: loc.Latitude, Longitude: loc.Longitude}

	// get route info for vehicle's current route
	if (loc.RouteID == nil) {
		// vehicle isn't on route
		log.Debugf("%s not on route", vehicle.Name)
		return
	}
	routeID := *loc.RouteID
	route, err := em.ms.Route(routeID)
	if err != nil {
		log.WithError(err).Error("unable to get route")
		return
	}

	lastDepartureTrack, err := em.getLastDepartureTrack(vehicle, route)
	if err != nil {
		log.WithError(err).Warn("unable to get last departure track")
		return
	}
	if len(lastDepartureTrack) < 1 {
		return
	}

	// is this track on the route?
	if !em.trackNearRoute(lastDepartureTrack, route) {
		log.Debug("track not near route")
		return
	}

	// stops visited in correct order?
	// inOrder, err := em.stopsVisitedInOrder(route, lastDepartureTrack)
	// if err != nil {
	// 	log.WithError(err).Error("unable to determine if stops visited in order")
	// 	return
	// }
	// if !inOrder {
	// 	log.Debug("stops not visited in order")
	// 	return
	// }

	durs, err := em.determineAverageTravelTimes(route)
	if err != nil {
		log.WithError(err).Error("unable to determine average travel times")
		return
	}

	// find index of stop zone on the route the vehicle is nearest to
	locIndices, err := em.determineNearestStopIndicesAlongRoute(lastDepartureTrack, route)
	if err != nil {
		log.WithError(err).Warn("unable to determine where vehicle is along route")
		return
	}
	locIndex := locIndices[len(locIndices) - 1]

	eta := &VehicleETA{
		VehicleID: vehicleID,
		RouteID: route.ID,
		StopETAs: make([]StopETA, 0, len(route.StopIDs)),
		Updated: time.Now(),
	}

	// stopPoints := make([]shuttletracker.Point, len(route.StopIDs))
	for i, stopID := range route.StopIDs {
		// find which zoneIndex this stop has
		zoneIdx := i

		// how many zones do we have to traverse to get there?
		traversal := zoneIdx - locIndex
		if traversal < 0 {
			// we passed the zone
			continue
		}

		// If this is zero, then the latest location is in the same zone as the stop.
		// This is useful to know since ETAs within a zone are probably not great.
		// Clients can just display a message about a vehicle arriving instead of
		// an ETA with a specific time.
		arriving := traversal == 0

		// add up zone travel durations
		totalDuration := time.Duration(0)
		for j := locIndex; j < zoneIdx; j++ {
			totalDuration += durs[j]
		}
		// last zone duration is half since stop is halfway through zone
		totalDuration += durs[zoneIdx]/2

		etaTime := loc.Created.Add(totalDuration)

		// sanity check
		if !etaTime.After(time.Now()) {
			log.Warn("ETA is in the past")
			continue
		}

		// would this ETA mean that the vehicle has to travel more than 35 mph (~15.6 meters/sec)?
		stop, err := em.ms.Stop(stopID)
		if err != nil {
			log.WithError(err).Error("unable to get stop")
			continue
		}
		stopPoint := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}
		directDistance := distanceBetween(locPoint, stopPoint)
		if directDistance / totalDuration.Seconds() > 15.6 {
			log.Warn("ETA is impossibly soon")
			continue
		}

		stopETA := StopETA{
			StopID: stopID,
			ETA: etaTime,
			Arriving: arriving,
		}
		eta.StopETAs = append(eta.StopETAs, stopETA)
	}

	log.Debugf("Calculated ETAs for %s", vehicle.Name)
	em.etaChan <- eta
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
