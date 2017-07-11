package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"

	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/cas.v1"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

// Configuration holds the settings for connecting to outside resources.
type Config struct {
	GoogleMapAPIKey      string
	GoogleMapMinDistance int
	CasURL               string `env:"CAS_URL"`
	Authenticate         bool   `env:"AUTHENTICATE" envDefault:"true"`
	ListenURL            string
}

// App holds references to Mongo resources.
type API struct {
	cfg     Config
	CasAUTH *cas.Client
	CasMEM  *cas.MemoryStore
	db      database.Database
	handler http.Handler
}

// InitApp initializes the application given a config and connects to backends.
// It also seeds any needed information to the database.
func New(cfg Config, db database.Database) (*API, error) {
	// Set up CAS authentication
	url, err := url.Parse(cfg.CasURL)
	if err != nil {
		return nil, err
	}
	var tickets *cas.MemoryStore

	client := cas.NewClient(&cas.Options{
		URL:   url,
		Store: nil,
	})

	// Create API instance to store database session and collections
	api := API{
		cfg:     cfg,
		CasAUTH: client,
		CasMEM:  tickets,
		db:      db,
	}

	// Read vehicle configuration file
	/*serr := readSeedConfiguration("seed/vehicle_seed.json", &app)
	if serr != nil {
		log.Fatalf("error reading vehicle configuration file: %v", serr)
	}*/

	r := mux.NewRouter()

	// Public
	r.HandleFunc("/vehicles", api.VehiclesHandler).Methods("GET")
	r.HandleFunc("/updates", api.UpdatesHandler).Methods("GET")
	r.HandleFunc("/updates/message", api.UpdateMessageHandler).Methods("GET")
	r.HandleFunc("/routes", api.RoutesHandler).Methods("GET")
	r.HandleFunc("/stops", api.StopsHandler).Methods("GET")

	// Admin
	r.Handle("/admin/", api.CasAUTH.HandleFunc(api.AdminHandler)).Methods("GET")
	r.Handle("/admin", api.CasAUTH.HandleFunc(api.AdminHandler)).Methods("GET")
	r.Handle("/admin/success/", api.CasAUTH.HandleFunc(api.AdminPageServer)).Methods("GET")
	r.Handle("/admin/success", api.CasAUTH.HandleFunc(api.AdminPageServer)).Methods("GET")
	r.Handle("/admin/logout/", api.CasAUTH.HandleFunc(api.AdminLogout)).Methods("GET")
	r.Handle("/admin/logout", api.CasAUTH.HandleFunc(api.AdminLogout)).Methods("GET")
	r.Handle("/vehicles/create", api.CasAUTH.HandleFunc(api.VehiclesCreateHandler)).Methods("POST")
	r.Handle("/vehicles/edit", api.CasAUTH.HandleFunc(api.VehiclesEditHandler)).Methods("POST")
	r.Handle("/vehicles/{id:[0-9]+}", api.CasAUTH.HandleFunc(api.VehiclesDeleteHandler)).Methods("DELETE")
	r.Handle("/routes/create", api.CasAUTH.HandleFunc(api.RoutesCreateHandler)).Methods("POST")
	r.Handle("/routes/{id:.+}", api.CasAUTH.HandleFunc(api.RoutesDeleteHandler)).Methods("DELETE")
	r.Handle("/stops/create", api.CasAUTH.HandleFunc(api.StopsCreateHandler)).Methods("POST")
	r.Handle("/stops/{id:.+}", api.CasAUTH.HandleFunc(api.StopsDeleteHandler)).Methods("DELETE")
	//r.HandleFunc("/import", api.ImportHandler).Methods("GET")

	// Legacy routes to support the ancient iOS app
	r.HandleFunc("/vehicles/current.js", api.LegacyVehiclesHandler).Methods("GET")
	r.HandleFunc("/displays/netlink.js", api.LegacyRoutesHandler).Methods("GET")

	// Static files
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Serve requests
	hand := api.CasAUTH.Handle(r)
	api.handler = hand

	return &api, nil
}

func NewConfig() *Config {
	return &Config{ListenURL: "localhost:8080"}
}

func (api *API) Run() {
	if err := http.ListenAndServe(api.cfg.ListenURL, api.handler); err != nil {
		log.WithError(err).Error("Unable to serve.")
	}
}

// IndexHandler serves the index page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

//readSeedConfiguration adds a new vehicle to the database from seed.
/*func readSeedConfiguration(fileName string, app *App) error {
	// Open seed_vehicle config file and decode JSON to app struct
	file, err := os.Open(fileName)

	// Error handling
	if err != nil {
		log.Warn(err)
	}
	// Create a decoder for a file
	fileread := json.NewDecoder(file)

	// Create map for json data and slice for vehicles
	var vehiclesMap map[string][]map[string]interface{} // map with string as key and ,list of map with string as key and anything as value, as value
	Vehicles := []Vehicle{}                             // list of default vehicle object

	// Call decode on fileread to place items into map
	if err := fileread.Decode(&vehiclesMap); err != nil {
		log.Warn(err)
	}

	// Initialize our vehicles
	for i := range vehiclesMap["Vehicles"] {
		item := vehiclesMap["Vehicles"][i]
		VehicleID, _ := item["VehicleID"].(string)
		VehicleName, _ := item["VehicleName"].(string)
		vehicle := Vehicle{VehicleID, VehicleName, time.Now(), time.Now(),false}
		Vehicles = append(Vehicles, vehicle)
	}

	// Add vehicles to the database
	for j := range Vehicles {
		api.Vehicles.Insert(&Vehicles[j])
	}

	return nil
}*/

type User struct {
	Name string
}

// AdminHandler serves the admin page.
func (api *API) AdminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%u", api.cfg.Authenticate)
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	} else {
		valid := false
		var users []User
		api.db.Users.Find(bson.M{}).All(&users)
		for _, u := range users {
			if u.Name == strings.ToLower(cas.Username(r)) {
				valid = true
			}
		}
		if api.cfg.Authenticate == false {
			valid = true
			fmt.Printf("not authenticating")
		}
		if valid {
			http.Redirect(w, r, "/admin/success/", 301)
		} else {
			http.Redirect(w, r, "/admin/logout/", 301)
		}
	}

}

func (api *API) AdminPageServer(w http.ResponseWriter, r *http.Request) {

	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		http.Redirect(w, r, "/admin/", 301)
		return
	} else {
		http.ServeFile(w, r, "admin.html")
	}

}

func (api *API) AdminLogout(w http.ResponseWriter, r *http.Request) {

	if cas.IsAuthenticated(r) {
		cas.RedirectToLogout(w, r)
	}

}

// WriteJSON writes the data as JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}
