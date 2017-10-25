package updater

import (
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"
)

type Updater struct {
	cfg            Config
	updateInterval time.Duration
	db             database.Database
	dataRegexp     *regexp.Regexp
}

type Config struct {
	DataFeed       string
	UpdateInterval string
}

func New(cfg Config, db database.Database) (*Updater, error) {
	updater := &Updater{cfg: cfg, db: db}

	interval, err := time.ParseDuration(cfg.UpdateInterval)
	if err != nil {
		return nil, err
	}
	updater.updateInterval = interval

	// Match each API field with any number (+)
	//   of the previous expressions (\d digit, \. escaped period, - negative number)
	//   Specify named capturing groups to store each field from data feed
	updater.dataRegexp = regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)

	return updater, nil
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		UpdateInterval: "10s",
	}
	v.SetDefault("updater.updateinterval", cfg.UpdateInterval)
	v.SetDefault("updater.datafeed", cfg.DataFeed)
	return cfg
}

// Run updater forever.
func (u *Updater) Run() {
	log.Debug("Updater started.")
	ticker := time.Tick(u.updateInterval)

	// Do one initial update.
	u.update()

	// Call update() every updateInterval.
	for range ticker {
		u.update()
	}
}

// Send a request to iTrak API, get updated shuttle info,
// store updated records in the database, and remove old records.
func (u *Updater) update() {
	// Make request to iTrak data feed
	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(u.cfg.DataFeed)
	if err != nil {
		log.WithError(err).Error("Could not get data feed.")
		return
	}

	// Read response body content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Could not read data feed.")
		return
	}
	resp.Body.Close()

	delim := "eof"
	// split the body of response by delimiter
	vehiclesData := strings.Split(string(body), delim)
	vehiclesData = vehiclesData[:len(vehiclesData)-1] // last element is EOF

	// TODO: Figure out if this handles == 1 vehicle correctly or always assumes > 1.
	if len(vehiclesData) <= 1 {
		log.Warnf("Found no vehicles delineated by '%s'.", delim)
	}

	wg := sync.WaitGroup{}
	// for parsed data, update each vehicle
	for _, vehicleData := range vehiclesData {
		wg.Add(1)
		go func(vehicleData string) {
			defer wg.Done()
			match := u.dataRegexp.FindAllStringSubmatch(vehicleData, -1)[0]
			// Store named capturing group and matching expression as a key value pair
			result := map[string]string{}
			for i, item := range match {
				result[u.dataRegexp.SubexpNames()[i]] = item
			}

			// Create new vehicle update & insert update into database

			// convert KPH to MPH
			speedKMH, err := strconv.ParseFloat(strings.Replace(result["speed"], "spd:", "", -1), 64)
			if err != nil {
				log.Error(err)
				return
			}
			speedMPH := kphToMPH(speedKMH)
			speedMPHString := strconv.FormatFloat(speedMPH, 'f', 5, 64)

			vehicle := model.Vehicle{}
			route := model.Route{}

			vehicleID := strings.Replace(result["id"], "Vehicle ID:", "", -1)
			err = u.db.Vehicles.Find(bson.M{"vehicleID": vehicleID}).One(&vehicle)
			if err == mgo.ErrNotFound {
				log.Warnf("Unknown vehicle ID \"%s\" returned by iTrak. Make sure all vehicles have been added.", vehicleID)
			} else if err != nil {
				log.WithError(err).Error("Unable to fetch vehicle.")
				return
			} else {
				// vehicle found and no error
				route, err = u.GuessRouteForVehicle(&vehicle)
				if err != nil {
					log.WithError(err).Error("Unable to guess route for vehicle.")
					return
				}
			}

			update := model.VehicleUpdate{
				VehicleID: strings.Replace(result["id"], "Vehicle ID:", "", -1),
				Lat:       strings.Replace(result["lat"], "lat:", "", -1),
				Lng:       strings.Replace(result["lng"], "lon:", "", -1),
				Heading:   strings.Replace(result["heading"], "dir:", "", -1),
				Speed:     speedMPHString,
				Lock:      strings.Replace(result["lock"], "lck:", "", -1),
				Time:      strings.Replace(result["time"], "time:", "", -1),
				Date:      strings.Replace(result["date"], "date:", "", -1),
				Status:    strings.Replace(result["status"], "trig:", "", -1),
				Created:   time.Now(),
				Route:     route.ID}

			if err := u.db.Updates.Insert(&update); err != nil {
				log.WithError(err).Errorf("Could not insert vehicle update.")
			}
		}(vehicleData)
	}
	wg.Wait()
	log.Debugf("Updated %d vehicles.", len(vehiclesData))

	// Prune updates older than one month
	info, err := u.db.Updates.RemoveAll(bson.M{"created": bson.M{"$lt": time.Now().AddDate(0, -1, 0)}})
	if err != nil {
		log.WithError(err).Error("Unable to remove old updates.")
		return
	}
	log.Debugf("Removed %d old updates.", info.Removed)
}

// Convert kmh to mph
func kphToMPH(kmh float64) float64 {
	return kmh * 0.621371192
}

//LastUpdatesForVehicle returns the last n updates for a vehicle
func (u *Updater) LastUpdatesForVehicle(vehicle *model.Vehicle, count int) (updates []model.VehicleUpdate) {
	err := u.db.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID}).Sort("-created").Limit(count).All(&updates)
	if err != nil {
		log.Error(err)
	}
	return
}

// GuessRouteForVehicle returns a guess at what route the vehicle is on.
// It may return an empty route if it does not believe a vehicle is on any route.
func (u *Updater) GuessRouteForVehicle(vehicle *model.Vehicle) (route model.Route, err error) {
	samples := 100
	var routes []model.Route
	err = u.db.Routes.Find(bson.M{}).All(&routes)
	if err != nil {
		log.Error(err)
	}

	routeDistances := make(map[string]float64)
	for _, route := range routes {
		routeDistances[route.ID] = 0
	}

	updates := u.LastUpdatesForVehicle(vehicle, samples)

	for _, update := range updates {
		updateLatitude, err := strconv.ParseFloat(update.Lat, 64)
		if err != nil {
			log.Error(err)
		}
		updateLongitude, err := strconv.ParseFloat(update.Lng, 64)
		if err != nil {
			log.Error(err)
		}

		for _, route := range routes {
			if !route.Enabled {
				routeDistances[route.ID] += math.Inf(0)
			}
			nearestDistance := math.Inf(0)
			for _, coord := range route.Coords {
				distance := math.Sqrt(math.Pow(updateLatitude-coord.Lat, 2) +
					math.Pow(updateLongitude-coord.Lng, 2))
				if distance < nearestDistance {
					nearestDistance = distance

				}
			}
			if nearestDistance > .003 {
				nearestDistance += 50
			}
			routeDistances[route.ID] += nearestDistance
		}
	}

	minDistance := math.Inf(0)
	var minRouteID string
	for id := range routeDistances {
		distance := routeDistances[id] / float64(samples)
		if distance < minDistance {
			minDistance = distance
			minRouteID = id
			// If more than ~5% of the last 100 samples were far away from a route, say the shuttle is not on a route
			// This is extremely aggressive and requires a shuttle to be on a route for ~5 minutes before it registers as on the route
			if minDistance > 5 {
				minRouteID = ""
			}
			log.Debugf("%v distance from nearest route: %v\n", vehicle.VehicleName, minDistance)

		}
	}

	// not on a route
	if minRouteID == "" {
		return model.Route{}, nil
	}

	err = u.db.Routes.Find(bson.M{"id": minRouteID}).One(&route)
	return route, err
}
