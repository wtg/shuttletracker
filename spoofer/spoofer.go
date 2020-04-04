package spoofer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)

// Updater handles periodically grabbing the latest vehicle location data from iTrak.
type Spoofer struct {
	cfg           Config
	spoofInterval time.Duration
	SpoofUpdates  bool
	spoofIndex    int
	dataRegexp    *regexp.Regexp
	ms            shuttletracker.ModelService
	mutex         *sync.Mutex
	sm            *sync.Mutex
	subscribers   []func(*shuttletracker.Location)
}

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

	// Match each API field with any number (+)
	//   of the previous expressions (\d digit, \. escaped period, - negative number)
	//   Specify named capturing groups to store each field from data feed
	spoofer.dataRegexp = regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)

	spoofer.spoofIndex = 0

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

		// Do one initial spoof.
		s.spoof()

		// Spoof updates every spoofInterval
		for range ticker {
			s.spoof()
		}
	}
}

// Subscribe allows callers to provide a function that is called after Updater parses a new Location.
func (s *Spoofer) Subscribe(f func(*shuttletracker.Location)) {
	s.sm.Lock()
	s.subscribers = append(s.subscribers, f)
	s.sm.Unlock()
}

func (s *Spoofer) notifySubscribers(loc *shuttletracker.Location) {
	s.sm.Lock()
	for _, sub := range s.subscribers {
		go sub(loc)
	}
	s.sm.Unlock()
}

// Spoofs locations for each shuttle
func (s *Spoofer) spoof() {
	wd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Errorf("Error getting working directory")
		return
	}
	filenamePrefix := wd + "/spoof_data/update" + strconv.Itoa(s.spoofIndex)
	vehicleId := 0
	filename := filenamePrefix + "-vehicle" + strconv.Itoa(vehicleId) + ".json"
	_, err = os.Stat(filename)
	for err == nil {
		spooffile, err := os.Open(filename)
		if err != nil {
			log.WithError(err).Errorf("Error opening update %d file", s.spoofIndex)
			return
		}
		bytes, err := ioutil.ReadAll(spooffile)
		if err != nil {
			log.WithError(err).Errorf("Error spoofing update %d", s.spoofIndex)
		}
		var update shuttletracker.Location
		json.Unmarshal(bytes, &update)

		update.Created = time.Now()
		update.Time = time.Now()

		if err := s.ms.CreateLocation(&update); err != nil {
			log.WithError(err).Errorf("Could not create spoofed location")
			return
		}
		log.Debugf("Spoofed location for vehicle %d", vehicleId)

		s.notifySubscribers(&update)

		vehicleId += 1
		filename = filenamePrefix + "-vehicle" + strconv.Itoa(vehicleId) + ".json"
		if _, err = os.Stat(filename); err != nil {
			break
		}
	}

	s.spoofIndex += 1
	spoofIndexTestFilename := wd + "/spoof_data/update" + strconv.Itoa(s.spoofIndex) + "-vehicle0.json"
	if _, err = os.Stat(spoofIndexTestFilename); err != nil {
		s.spoofIndex = 0
	}
}
