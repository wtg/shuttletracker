package smooth

import (
	"math"
	"sync"
	"time"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/updater"
)

type SmoothTrackingManager struct {
	cfg                Config
	ms                 shuttletracker.ModelService
	predictions        map[int64]*Prediction
	predictionInterval time.Duration
	predictUpdates     bool
	updates            map[int64]*shuttletracker.Location
	vehicleIDs         []int64
	sm                 *sync.Mutex
	subscribers        []func(Prediction)
}

type Config struct {
	PredictionInterval string
	PredictUpdates     bool
}

func NewManager(cfg Config, ms shuttletracker.ModelService, updater *updater.Updater) (*SmoothTrackingManager, error) {
	stm := &SmoothTrackingManager{
		cfg:         cfg,
		ms:          ms,
		predictions: map[int64]*Prediction{},
		updates:     map[int64]*shuttletracker.Location{},
		sm:          &sync.Mutex{},
		subscribers: []func(Prediction){},
	}

	interval, err := time.ParseDuration(cfg.PredictionInterval)
	if err != nil {
		return nil, err
	}
	stm.predictionInterval = interval
	stm.predictUpdates = cfg.PredictUpdates

	// Subscribe to new Locations with Updater
	updater.Subscribe(stm.locationSubscriber)

	return stm, nil
}

func (stm *SmoothTrackingManager) locationSubscriber(loc *shuttletracker.Location) {
	if loc.VehicleID == nil {
		return
	}
	stm.updates[*loc.VehicleID] = loc
	index := -1
	for i, id := range stm.vehicleIDs {
		if id == *loc.VehicleID {
			index = i
			break
		}
	}

	if loc.RouteID != nil {
		if index < 0 {
			stm.vehicleIDs = append(stm.vehicleIDs, *loc.VehicleID)
		}
	} else if index >= 0 {
		// This vehicle is no longer on a route; remove it from the active vehicles list
		stm.vehicleIDs[index] = stm.vehicleIDs[len(stm.vehicleIDs)-1]
		stm.vehicleIDs[len(stm.vehicleIDs)-1] = 0
		stm.vehicleIDs = stm.vehicleIDs[:len(stm.vehicleIDs)-1]
	}

	stm.sm.Lock()
	prediction, exists := stm.predictions[*loc.VehicleID]
	stm.sm.Unlock()

	if exists {
		diffIndex := int64(math.Abs(float64(prediction.Index - index)))
		diffDistance := DistanceBetween(prediction.Point, shuttletracker.Point{Latitude: loc.Latitude, Longitude: loc.Longitude})
		log.Debugf("UPDATED VEHICLE %d", *loc.VehicleID)
		log.Debugf("Predicted: %d, (%f, %f)", prediction.Index, prediction.Point.Latitude, prediction.Point.Longitude)
		log.Debugf("Actual: %d, (%f, %f)", index, loc.Latitude, loc.Longitude)
		log.Debugf("Difference: %d points or %f meters", diffIndex, diffDistance)
	}
}

func (stm *SmoothTrackingManager) predict() {
	wg := sync.WaitGroup{}
	for _, id := range stm.vehicleIDs {
		wg.Add(1)
		go func(id int64) {
			stm.predictVehiclePosition(id)
		}(id)
	}
	wg.Wait()
}

func (stm *SmoothTrackingManager) predictVehiclePosition(vehicleID int64) {
	vehicle, err := stm.ms.Vehicle(vehicleID)
	if err != nil {
		log.WithError(err).Errorf("cannot get vehicle for prediction")
	}
	update, exists := stm.updates[vehicle.ID]
	if !exists {
		log.Errorf("no prior update for vehicle")
	}
	if update.RouteID == nil {
		log.Errorf("no route for vehicle")
	}
	route, err := stm.ms.Route(*update.RouteID)
	if err != nil {
		log.WithError(err).Errorf("cannot get route for prediction")
	}
	prediction := NaivePredictPosition(vehicle, update, route)
	newUpdate := &shuttletracker.Location{
		TrackerID: update.TrackerID,
		Latitude:  prediction.Point.Latitude,
		Longitude: prediction.Point.Longitude,
		Heading:   update.Heading,
		Speed:     update.Speed,
		Time:      time.Now(),
		RouteID:   &route.ID,
	}
	if err := stm.ms.CreateLocation(newUpdate); err != nil {
		log.WithError(err).Error("could not create location for prediction")
	}
	stm.handleNewPrediction(&prediction)
}

// Run is in charge of managing all of the state inside of ETAManager.
func (stm *SmoothTrackingManager) Run() {
	if stm.predictUpdates {
		ticker := time.Tick(stm.predictionInterval)
		for range ticker {
			stm.predict()
		}
	}
}

func (stm *SmoothTrackingManager) handleNewPrediction(prediction *Prediction) {
	stm.sm.Lock()
	stm.predictions[prediction.VehicleID] = prediction
	for _, sub := range stm.subscribers {
		sub(*prediction)
	}
	stm.sm.Unlock()
}

// Subscribe allows callers to provide a callback to receive new VehicleETAs.
func (stm *SmoothTrackingManager) Subscribe(sub func(Prediction)) {
	stm.sm.Lock()
	stm.subscribers = append(stm.subscribers, sub)
	stm.sm.Unlock()
}
