package updater

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"math"

	"github.com/spf13/viper"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"
	"gopkg.in/mgo.v2/bson"

)

var (
	// Match each API field with any number (+)
	//   of the previous expressions (\d digit, \. escaped period, - negative number)
	//   Specify named capturing groups to store each field from data feed
	dataRe     = regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)
	dataNames  = dataRe.SubexpNames()
	lastUpdate time.Time
)

type Updater struct {
	cfg            Config
	updateInterval time.Duration
	db             database.Database
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

// Send a request to iTrak API, get updated shuttle info, and
// finally store updated records in the database.
func (u *Updater) update() {
	// Make request to our tracking data feed
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
	// BUG: if the request fails, it will give undefined result

	// TODO: Figure out if this handles == 1 vehicle correctly or always assumes > 1.
	if len(vehiclesData) <= 1 {
		log.Warnf("found no vehicles delineated by '%s'", delim)
	}

	updated := 0
	// for parsed data, update each vehicle
	for i := 0; i < len(vehiclesData)-1; i++ {
		match := dataRe.FindAllStringSubmatch(vehiclesData[i], -1)[0]
		// Store named capturing group and matching expression as a key value pair
		result := map[string]string{}
		for i, item := range match {
			result[dataNames[i]] = item
		}

		// Create new vehicle update & insert update into database
		// add computation of segment that the shuttle resides on and the arrival time to next N stops [here]

		// convert KPH to MPH
		speedKMH, err := strconv.ParseFloat(strings.Replace(result["speed"], "spd:", "", -1), 64)
		if err != nil {
			log.Error(err)
			continue
		}
		speedMPH := kphToMPH(speedKMH)
		speedMPHString := strconv.FormatFloat(speedMPH, 'f', 5, 64)
		vehicle := model.Vehicle{}
		route := model.Route{}

		error := u.db.Vehicles.Find(bson.M{"vehicleID": strings.Replace(result["id"], "Vehicle ID:", "", -1)}).One(&vehicle)
		if (error == nil){
			 route = u.GuessRouteForVehicle(&vehicle)
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
			Route:		 route.ID}

		// convert updated time to local time
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			log.WithError(err).Error("Could not load time zone information.")
			continue
		}
		lastUpdate = time.Now().In(loc)

		if err := u.db.Updates.Insert(&update); err != nil {
			log.WithError(err).Errorf("Could not insert vehicle update.")
		} else {
			updated++
		}

		// here if parsing error, updated will be incremented, wait, the whole thing will crash, isn't it?
	}
	log.Debugf("Successfully updated %d/%d vehicles.", updated, len(vehiclesData)-1)
}

// Convert kmh to mph
func kphToMPH(kmh float64) float64 {
	return kmh * 0.621371192
}


func (api *Updater) LastUpdatesForVehicle(vehicle *model.Vehicle, count int) (updates []model.VehicleUpdate) {
	err := api.db.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID}).Sort("-created").Limit(count).All(&updates)
	if err != nil {
		log.Error(err)
	}
	return
}

func (api *Updater) GuessRouteForVehicle(vehicle *model.Vehicle) (route model.Route) {
	samples := 100
	var routes []model.Route
	err := api.db.Routes.Find(bson.M{}).All(&routes)
	if err != nil {
		log.Error(err)
	}

	routeDistances := make(map[string]float64)
	for _, route := range routes {
		routeDistances[route.ID] = 0
	}

	updates := api.LastUpdatesForVehicle(vehicle, samples)

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
			nearestDistance := math.Inf(0)
			for _, coord := range route.Coords {
				distance := math.Sqrt(math.Pow(updateLatitude - coord.Lat, 2) +
					math.Pow(updateLongitude - coord.Lng, 2))
				if distance < nearestDistance {
					nearestDistance = distance
				}
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
		}
	}

	err = api.db.Routes.Find(bson.M{"id": minRouteID}).One(&route)
	if err != nil {
		log.Error(err)
	}

	return
}
