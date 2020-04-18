package spoofer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)

// Spoofer creates "fake" (spoofed) updates at regular intervals from data in the spoof_data folder
type Spoofer struct {
	cfg           Config
	spoofInterval time.Duration
	SpoofUpdates  bool
	spoofIndexes  map[int]int
	updates       map[int][]shuttletracker.Location
	updateID      int64
	ms            shuttletracker.ModelService
	mutex         *sync.Mutex
	sm            *sync.Mutex
	subscribers   []func(*shuttletracker.Location)
}

// Configuration for Spoofer; determines whether or not spoofed updates will be created and the interval
// at which it will occur
type Config struct {
	SpoofInterval string
	SpoofUpdates  bool
}

// Creates a new Spoofer
func New(cfg Config, ms shuttletracker.ModelService) (*Spoofer, error) {
	spoofer := &Spoofer{
		cfg:         cfg,
		ms:          ms,
		mutex:       &sync.Mutex{},
		sm:          &sync.Mutex{},
		subscribers: []func(*shuttletracker.Location){},
	}

	interval, err := time.ParseDuration(cfg.SpoofInterval)
	if err != nil {
		return nil, err
	}
	spoofer.spoofInterval = interval
	spoofer.SpoofUpdates = cfg.SpoofUpdates
	spoofer.spoofIndexes = make(map[int]int)
	spoofer.updates = make(map[int][]shuttletracker.Location)
	spoofer.updateID = 1

	return spoofer, nil
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		SpoofInterval: "10s",
		SpoofUpdates:  false,
	}
	v.SetDefault("spoof.spoofInterval", cfg.SpoofInterval)
	v.SetDefault("spoof.spoofUpdates", cfg.SpoofUpdates)
	return cfg
}

// Run spoofer forever.
func (s *Spoofer) Run() {
	if s.SpoofUpdates {
		log.Debug("Spoofer started.")
		ticker := time.Tick(s.spoofInterval)

		// Parse all update data
		s.parseUpdates()

		// Do one initial spoof
		s.spoof()

		// Spoof updates for each vehicle every spoofInterval
		for range ticker {
			s.spoof()
		}
	}
}

// Sequentially reads and caches all JSON data to create updates from
func (s *Spoofer) parseUpdates() {
	wd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Errorf("Error getting working directory")
		return
	}
	vehicleIndex := 0
	filename := wd + "/spoof_data/vehicle" + strconv.Itoa(vehicleIndex) + ".json"
	_, err = os.Stat(filename)
	for err == nil {
		vehiclefile, err := os.Open(filename)
		if err != nil {
			log.WithError(err).Errorf("Error opening vehicle %d file", vehicleIndex)
			return
		}
		bytes, err := ioutil.ReadAll(vehiclefile)
		if err != nil {
			log.WithError(err).Errorf("Error reading vehicle %d file", vehicleIndex)
		}
		var updates []shuttletracker.Location
		json.Unmarshal(bytes, &updates)

		log.Debugf("Read %d updates for vehicle %d", len(updates), vehicleIndex)

		// Only cache data for this vehicle if it has updates
		if len(updates) > 0 {
			s.updates[vehicleIndex] = updates
			s.spoofIndexes[vehicleIndex] = 0
		}

		vehicleIndex += 1
		filename = wd + "/spoof_data/vehicle" + strconv.Itoa(vehicleIndex) + ".json"
		if _, err = os.Stat(filename); err != nil {
			break
		}
	}
}

// Spoofs the next location for each vehicle
func (s *Spoofer) spoof() {
	for vehicleIndex, updates := range s.updates {
		update := updates[s.spoofIndexes[vehicleIndex]]
		update.Created = time.Now()
		update.Time = time.Now()
		update.TrackerID = "4572001148" // TODO
		update.ID = s.updateID
		if err := s.ms.CreateLocation(&update); err != nil {
			log.WithError(err).Errorf("Could not create spoofed location")
			return
		}
		log.Debugf("Spoofed location for vehicle %d", vehicleIndex)

		s.notifySubscribers(&update)

		s.spoofIndexes[vehicleIndex] += 1
		if s.spoofIndexes[vehicleIndex] >= len(updates) {
			s.spoofIndexes[vehicleIndex] = 0
		}
		s.updateID += 1
	}
}

// Subscribe allows callers to provide a function that is called after Spoofer creates a new location.
// Automatically reroutes all of Updater's subscribers to Spoofer when update spoofing is enabled.
func (s *Spoofer) Subscribe(f func(*shuttletracker.Location)) {
	s.sm.Lock()
	s.subscribers = append(s.subscribers, f)
	s.sm.Unlock()
}

// Notifies all of Spoofer's subscribers with a new location.
func (s *Spoofer) notifySubscribers(loc *shuttletracker.Location) {
	s.sm.Lock()
	for _, sub := range s.subscribers {
		go sub(loc)
	}
	s.sm.Unlock()
}
