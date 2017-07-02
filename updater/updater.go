package updater

import (
	"time"
	"net/http"

	"github.com/wtg/shuttletracker/log"
	"io/ioutil"
	"strings"
	"regexp"
	"strconv"
	"github.com/wtg/shuttletracker/api"
	"github.com/wtg/shuttletracker/database"
)

var (
	// Match each API field with any number (+)
	//   of the previous expressions (\d digit, \. escaped period, - negative number)
	//   Specify named capturing groups to store each field from data feed
	dataRe    = regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)
	dataNames = dataRe.SubexpNames()
	lastUpdate time.Time
)

type Updater struct {
	cfg Config
	updateInterval time.Duration
	db database.Database
}

type Config struct {
	DataFeed string
	UpdateInterval string
}

func New(cfg Config, db database.Database) (*Updater) {
	updater := &Updater{cfg: cfg, db: db}

	interval, err := time.ParseDuration(cfg.UpdateInterval)
	if err != nil {
		log.Error(err)
		return nil
	}
	updater.updateInterval = interval

	return updater
}

func NewConfig() *Config {
	return &Config{
		UpdateInterval: "10s",
	}

}

// Run updater forever.
func (u *Updater) Run() {
	for {
		u.update()
		time.Sleep(u.updateInterval)
	}
}

// UpdateShuttles send a request to iTrak API, gets updated shuttle info, and
// finally store updated records in db.
func (u *Updater) update() {
	var st time.Duration
	if st == 0 {
		// Initialize the sleep timer after the first sleep.  This lets us sleep during errors
		// when we 'continue' back to the top of the loop without waiting to sleep for the first
		// update run.
		st = time.Duration(u.updateInterval) * time.Second
	}

	// Make request to our tracking data feed
	resp, err := http.Get(u.cfg.DataFeed)
	if err != nil {
		log.Errorf("error getting data feed: %v", err)
		return
	}

	// Read response body content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading data feed: %v", err)
		return
	}

	// this cannot be deferred because that will never get garbage collected
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
		speedMPH := KPHtoMPH(speedKMH)
		speedMPHString := strconv.FormatFloat(speedMPH, 'f', 5, 64)

		update := api.VehicleUpdate{
			VehicleID: strings.Replace(result["id"], "Vehicle ID:", "", -1),
			Lat:       strings.Replace(result["lat"], "lat:", "", -1),
			Lng:       strings.Replace(result["lng"], "lon:", "", -1),
			Heading:   strings.Replace(result["heading"], "dir:", "", -1),
			Speed:     speedMPHString,
			Lock:      strings.Replace(result["lock"], "lck:", "", -1),
			Time:      strings.Replace(result["time"], "time:", "", -1),
			Date:      strings.Replace(result["date"], "date:", "", -1),
			Status:    strings.Replace(result["status"], "trig:", "", -1),
			Created:   time.Now()}

		// convert updated time to local time
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			log.Error(err.Error())
			continue
		}

		lastUpdate = time.Now().In(loc)

		if err := u.db.Updates.Insert(&update); err != nil {
			log.Errorf("error inserting vehicle update(%v): %v", update, err)
		} else {
			updated++
		}

		// here if parsing error, updated will be incremented, wait, the whole thing will crash, isn't it?
	}
	log.Debugf("Successfully updated %d/%d vehicles.", updated, len(vehiclesData)-1)
}

// convert kmh to mph
func KPHtoMPH(kmh float64) (mph float64) {
	mph = kmh * 0.621371192
	return
}
