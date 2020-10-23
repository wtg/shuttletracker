package smooth

import (
	"math"
	"sync"
	"time"

	"github.com/spf13/viper"
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
	debugMode          bool
	updates            map[int64]*shuttletracker.Location
	vehicleIDs         []int64
	sm                 *sync.Mutex
	subscribers        []func(Prediction)
	numDifferences     int
	averageDifference  float64
}

type Config struct {
	PredictionInterval string
	PredictUpdates     bool
	DebugMode          bool
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		PredictUpdates:     true,
		DebugMode:          true,
		PredictionInterval: "1s",
	}

	v.SetDefault("smooth.predictupdates", cfg.PredictUpdates)
	v.SetDefault("smooth.predictioninterval", cfg.PredictionInterval)
	v.SetDefault("smooth.debugmode", cfg.DebugMode)

	return cfg
}

// Creates a new SmoothTrackingManager
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
	stm.debugMode = cfg.DebugMode

	// Subscribe to new Locations with Updater
	updater.Subscribe(stm.locationSubscriber)

	return stm, nil
}

// If enabled, runs forever making predictions at regular intervals
func (stm *SmoothTrackingManager) Run() {
	if stm.predictUpdates {
		ticker := time.Tick(stm.predictionInterval)
		for range ticker {
			stm.predict()
		}
	}
}

// Asynchronously make predictions for all active vehicles
func (stm *SmoothTrackingManager) predict() {
	wg := sync.WaitGroup{}
	for _, id := range stm.vehicleIDs {
		wg.Add(1)
		go func(id int64) {
			stm.predictVehiclePosition(id)
			wg.Done()
		}(id)
	}
	wg.Wait()
}

// Use the prediction algorithm to make a prediction of this vehicle's current location, create a new location, and send to handleNewPrediction for further processing
func (stm *SmoothTrackingManager) predictVehiclePosition(vehicleID int64) {	
	vehicle, err := stm.ms.Vehicle(vehicleID)
	if err != nil {
		log.WithError(err).Errorf("Cannot get vehicle %d for prediction", vehicleID)
	}
	update, exists := stm.updates[vehicle.ID]
	if !exists {
		log.Errorf("No prior update for vehicle %d to base prediction on", vehicleID)
	}
	if update.RouteID == nil {
		log.Errorf("No route for vehicle %d to base prediction on", vehicleID)
	}
	route, err := stm.ms.Route(*update.RouteID)
	if err != nil {
		log.WithError(err).Errorf("Cannot get route for vehicle %d to base prediction on", vehicleID)
	}
	prediction := NaivePredictPosition(vehicle, update, route)
	newUpdate := &shuttletracker.Location{
		TrackerID: update.TrackerID,
		Latitude:  prediction.Point.Latitude,
		Longitude: prediction.Point.Longitude,
		Heading:   prediction.Angle,
		Speed:     update.Speed,
		Time:      time.Now(),
		RouteID:   &route.ID,
	}
	if err := stm.ms.CreateLocation(newUpdate); err != nil {
		log.WithError(err).Errorf("Cannot not create predicted location for vehicle %d", vehicleID)
	}
	stm.handleNewPrediction(&prediction)
}

// Put new prediction in the predictions map and notify subscribers
func (stm *SmoothTrackingManager) handleNewPrediction(prediction *Prediction) {
	stm.sm.Lock()
	stm.predictions[prediction.VehicleID] = prediction
	for _, sub := range stm.subscribers {
		sub(*prediction)
	}
	stm.sm.Unlock()
}

// Allows callers to provide a callback to receive new predictions when they are created
func (stm *SmoothTrackingManager) Subscribe(sub func(Prediction)) {
	stm.sm.Lock()
	stm.subscribers = append(stm.subscribers, sub)
	stm.sm.Unlock()
}

// Subscribe to new locations received from the updater
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

	// Output comparison between predicted and actual position
	if exists {
		diffIndex := int64(math.Abs(float64(prediction.Index - index)))
		diffDistance := DistanceBetween(prediction.Point, shuttletracker.Point{Latitude: loc.Latitude, Longitude: loc.Longitude})
		log.Debugf("UPDATED VEHICLE %d", *loc.VehicleID)
		log.Debugf("Predicted: %d, (%f, %f)", prediction.Index, prediction.Point.Latitude, prediction.Point.Longitude)
		log.Debugf("Actual: %d, (%f, %f)", index, loc.Latitude, loc.Longitude)
		log.Debugf("Difference: %d points or %f meters", diffIndex, diffDistance)

		// add statistics code to calculate stats
		if stm.debugMode {
			stm.numDifferences += 1
			stm.averageDifference = stm.averageDifference + (diffDistance - stm.averageDifference) / float64(stm.numDifferences)
			log.Debugf("Average Difference is %f", stm.averageDifference)
			log.Debugf("Number of Differences is %d", stm.numDifferences)
		}
	}

	if stm.debugMode {
		log.Debugf("Current number of predictions so far is %d", stm.numDifferences)
		log.Debugf("Average Difference is %f", stm.averageDifference)
		log.Debugf("Number of Differences is %d", stm.numDifferences)
	}
}
