package eta

import (
	"errors"
	"math"
	"sync"
	"time"

	// "github.com/wcharczuk/go-chart"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/updater"
)

// ETAManager implements ETAService and provides ETAs for Vehicles to Stops.
type ETAManager struct {
	ms          shuttletracker.ModelService
	etaChan     chan *shuttletracker.VehicleETA
	etas        map[int64]*shuttletracker.VehicleETA
	etasReqChan chan chan map[int64]shuttletracker.VehicleETA

	sm          *sync.Mutex
	subscribers []func(shuttletracker.VehicleETA)
}

// NewManager creates an ETAManager subscribed to Location updates from Updater.
func NewManager(ms shuttletracker.ModelService, updater *updater.Updater) (*ETAManager, error) {
	em := &ETAManager{
		ms:          ms,
		etaChan:     make(chan *shuttletracker.VehicleETA, 50),
		etas:        map[int64]*shuttletracker.VehicleETA{},
		etasReqChan: make(chan chan map[int64]shuttletracker.VehicleETA),
		sm:          &sync.Mutex{},
		subscribers: []func(shuttletracker.VehicleETA){},
	}

	// subscribe to new Locations with Updater
	updater.Subscribe(em.locationSubscriber)

	return em, nil
}

// This gets new Locations from Updater. As soon as this happens, we'll
// determine new ETAs for the vehicle in another goroutine.
func (em *ETAManager) locationSubscriber(loc *shuttletracker.Location) {
	go em.handleNewLocation(loc)
}

func (em *ETAManager) handleNewLocation(loc *shuttletracker.Location) {
	if loc.VehicleID == nil {
		// can't do anything...
		return
	}
	vehicleID := *loc.VehicleID
	eta, err := em.calculateVehicleETAs(vehicleID)
	if err != nil {
		log.WithError(err).Errorf("unable to calculate ETAs for vehicle ID %d", vehicleID)
		return
	}

	log.Debugf("calculated ETAs for vehicle ID %d", vehicleID)
	em.etaChan <- eta
}

func (em *ETAManager) calculateVehicleETAs(vehicleID int64) (*shuttletracker.VehicleETA, error) {
	// get vehicle info
	vehicle, err := em.ms.Vehicle(vehicleID)
	if err != nil {
		return nil, err
	}

	loc, err := em.ms.LatestLocation(vehicleID)
	if err != nil {
		return nil, err
	}

	locPoint := shuttletracker.Point{Latitude: loc.Latitude, Longitude: loc.Longitude}

	eta := &shuttletracker.VehicleETA{
		VehicleID: vehicleID,
		StopETAs:  []shuttletracker.StopETA{},
		Updated:   time.Now(),
	}

	// get route info for vehicle's current route
	if loc.RouteID == nil {
		// vehicle isn't on route
		return eta, nil
	}

	routeID := *loc.RouteID
	eta.RouteID = routeID
	route, err := em.ms.Route(routeID)
	if err != nil {
		return nil, err
	}

	lastDepartureTrack, err := em.getLastDepartureTrack(vehicle, route)
	if err != nil {
		return nil, err
	}
	if len(lastDepartureTrack) < 1 {
		return eta, nil
	}

	// is this track on the route?
	if !em.trackNearRoute(lastDepartureTrack, route) {
		return eta, nil
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
		return nil, err
	}

	// find index of stop zone on the route the vehicle is nearest to
	locIndices, err := em.determineNearestStopIndicesAlongRoute(lastDepartureTrack, route)
	if err != nil {
		return nil, err
	}

	// The last locIndex is where the vehicle is now. If the track is short, we can't determine this.
	if len(locIndices) == 0 {
		return eta, nil
	}
	locIndex := locIndices[len(locIndices)-1]

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
		totalDuration += durs[zoneIdx] / 2

		etaTime := loc.Created.Add(totalDuration)

		// sanity check
		if !etaTime.After(time.Now()) {
			log.Warn("ETA is in the past")
			continue
		}

		// would this ETA mean that the vehicle has to travel more than 35 mph (~15.6 meters/sec)?
		stop, err := em.ms.Stop(stopID)
		if err != nil {
			return nil, err
		}
		stopPoint := shuttletracker.Point{Latitude: stop.Latitude, Longitude: stop.Longitude}
		directDistance := distanceBetween(locPoint, stopPoint)
		if directDistance/totalDuration.Seconds() > 15.6 {
			log.Warn("ETA is impossibly soon")
			continue
		}

		stopETA := shuttletracker.StopETA{
			StopID:   stopID,
			ETA:      etaTime,
			Arriving: arriving,
		}
		eta.StopETAs = append(eta.StopETAs, stopETA)
	}

	return eta, nil
}

// Run is in charge of managing all of the state inside of ETAManager.
func (em *ETAManager) Run() {
	err := em.createInitialETAs()
	if err != nil {
		log.WithError(err).Error("unable to create initial ETAs")
	}

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

func (em *ETAManager) createInitialETAs() error {
	vehicles, err := em.ms.Vehicles()
	if err != nil {
		return err
	}
	for _, vehicle := range vehicles {
		eta, err := em.calculateVehicleETAs(vehicle.ID)
		if err != nil {
			log.WithError(err).Errorf("unable to calculate ETAs for vehicle ID %d", vehicle.ID)
			continue
		}
		log.Debugf("calculated ETAs for vehicle ID %d", vehicle.ID)
		em.etaChan <- eta
	}
	return nil
}

func (em *ETAManager) handleNewETA(eta *shuttletracker.VehicleETA) {
	em.etas[eta.VehicleID] = eta

	// notify subscribers
	em.sm.Lock()
	for _, sub := range em.subscribers {
		sub(*eta)
	}
	em.sm.Unlock()
}

// spit out all current ETAs over the provided channel
func (em *ETAManager) processETAsRequest(c chan map[int64]shuttletracker.VehicleETA) {
	etas := map[int64]shuttletracker.VehicleETA{}
	for k, v := range em.etas {
		etas[k] = *v
	}

	c <- etas
}

// Iterate over all ETAs and remove those that have expired.
// We also send empty ETAs after we clean them up.
func (em *ETAManager) cleanup() {
	log.Debug("ETAManager cleanup")
	now := time.Now()
	for _, vehicleETA := range em.etas {
		stopETAs := vehicleETA.StopETAs
		shouldPush := false
		for i := len(stopETAs) - 1; i >= 0; i-- {
			stopETA := stopETAs[i]
			if now.After(stopETA.ETA) {
				shouldPush = true
				stopETAs = append(stopETAs[:i], stopETAs[i+1:]...)
			}
		}
		vehicleETA.StopETAs = stopETAs
		if shouldPush {
			em.etaChan <- vehicleETA
		}
	}
}

// Subscribe allows callers to provide a callback to receive new VehicleETAs.
func (em *ETAManager) Subscribe(sub func(shuttletracker.VehicleETA)) {
	em.sm.Lock()
	em.subscribers = append(em.subscribers, sub)
	em.sm.Unlock()
}

// CurrentETAs can be called by anyone to get ETAManager's current view of vehicle ETAs.
// It returns structs as values in order to prevent data races.
func (em *ETAManager) CurrentETAs() map[int64]shuttletracker.VehicleETA {
	etasChan := make(chan map[int64]shuttletracker.VehicleETA)
	em.etasReqChan <- etasChan
	return <-etasChan
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
		if ld.dist > 30 {
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
			if ld.dist < 30 {
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
		if track[len(track)-1].loc.Time.Sub(track[0].loc.Time) < 5*time.Minute {
			continue
		}

		// total distance traveled is at least 75% of route length? no more than
		// 110% of route length? (it is unlikely that a track is longer than a route
		// since a route has many more points than a track—a track essentially cuts corners.)
		routeLength := calculateRouteDistance(route)
		d := calculateDistance(track)
		if d < routeLength*0.75 || d > routeLength*1.1 {
			continue
		}

		// no point on the track more than 100 m from the route?
		locTrack := make([]*shuttletracker.Location, len(track))
		for i, ld := range track {
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

// WARNING: this assumes that the initial stop only occurs on a route _once_!
func (em *ETAManager) getLastDepartureTrack(vehicle *shuttletracker.Vehicle, route *shuttletracker.Route) ([]*shuttletracker.Location, error) {
	since := time.Now().Add(time.Minute * -30)
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
		return []*shuttletracker.Location{}, nil
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
	for i := 0; i < len(locs)-1; i++ {
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
	if len(route.StopIDs)-len(desiredStopIDs) > 2 {
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

func (em *ETAManager) determineNearestStopIndicesAlongRoute(track []*shuttletracker.Location, route *shuttletracker.Route) ([]int, error) {
	if len(track) < 1 {
		return []int{}, nil
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
