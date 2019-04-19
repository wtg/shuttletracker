package api

import (
	"encoding/json"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"

	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)


const (
	vapidPublicKey = "BHu_01FAmOhIaQ1KXX4qqHiJ7ire9s5dYTK4TF2dFXbeWb0fFvfpjJl3zaQjonIjhx1bl7IlQ_MWFsQBzAYZV9I"
)
var vapidPrivateKey string

type Request struct {
	Delay 			int64 		`json:"delay"`
	Campus			string 		`json:"campus"`
	Subscription 	string      `json:"subscription"`
}

type Subscript struct {
	Subscription   	string     	`json:"subscript"`   
}

// Config holds API settings.
type Config struct {
	GoogleMapAPIKey      string
	GoogleMapMinDistance int
	CasURL               string
	Authenticate         bool
	ListenURL            string
	MapboxAPIKey         string
	PrivateVapidKey		 string
}

// API is responsible for configuring handlers for HTTP endpoints.
type API struct {
	cfg        Config
	handler    http.Handler
	ms         shuttletracker.ModelService
	msg        shuttletracker.MessageService
	updater    shuttletracker.UpdaterService
	fm         *fusionManager
	etaManager shuttletracker.ETAService
}

// New initializes the application given a config and connects to backends.
// It also seeds any needed information to the database.
func New(cfg Config, ms shuttletracker.ModelService, msg shuttletracker.MessageService, us shuttletracker.UserService, updater shuttletracker.UpdaterService, etaManager shuttletracker.ETAService) (*API, error) {
	// Set up CAS authentication
	url, err := url.Parse(cfg.CasURL)
	if err != nil {
		return nil, err
	}

	// Set up fusion manager
	fm, err := newFusionManager(etaManager, ms)
	if err != nil {
		return nil, err
	}

	// Create API instance to store database session and collections
	api := API{
		cfg:        cfg,
		ms:         ms,
		msg:        msg,
		updater:    updater,
		fm:         fm,
		etaManager: etaManager,
	}

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)
	r.Use(etag)

	cli := CreateCASClient(url, us, cfg.Authenticate)

	// Vehicles
	r.Route("/vehicles", func(r chi.Router) {
		r.Get("/", api.VehiclesHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/create", api.VehiclesCreateHandler)
			r.Post("/edit", api.VehiclesEditHandler)
			r.Delete("/", api.VehiclesDeleteHandler)
		})
	})

	// Updates
	r.Route("/updates", func(r chi.Router) {
		r.Get("/", api.UpdatesHandler)
	})

	// History
	r.Route("/history", func(r chi.Router) {
		r.Get("/", api.HistoryHandler)
	})

	// Admin message
	r.Route("/adminMessage", func(r chi.Router) {
		r.Get("/", api.AdminMessageHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/", api.SetAdminMessage)
		})
	})

	// Routes
	r.Route("/routes", func(r chi.Router) {
		r.Get("/", api.RoutesHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/create", api.RoutesCreateHandler)
			r.Post("/edit", api.RoutesEditHandler)
			r.Delete("/", api.RoutesDeleteHandler)
		})
	})

	// Stops
	r.Route("/stops", func(r chi.Router) {
		r.Get("/", api.StopsHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/create", api.StopsCreateHandler)
			r.Delete("/", api.StopsDeleteHandler)
		})
	})

	// Fusion
	r.Mount("/fusion", api.fm.router(cli.casauth))

	r.Get("/logout/", cli.logout)
	// Admin
	r.Route("/admin", func(r chi.Router) {
		r.Use(cli.casauth)
		r.Get("/*", api.AdminHandler)
		r.Get("/login", api.AdminHandler)
		r.Get("/logout", cli.logout)
	})

	r.Group(func(r chi.Router) {
		r.Use(cli.casauth)
		r.Get("/getKey/", api.KeyHandler)
	})

	r.Method("GET", "/static/*", http.StripPrefix("/static/", http.FileServer(staticFileSystem{http.Dir("static/")})))

	r.Get("/", api.IndexHandler)
	r.Get("/about", api.IndexHandler)
	r.Get("/schedules", api.IndexHandler)
	r.Get("/settings", api.IndexHandler)
	r.Get("/serviceworker.js", api.ServiceWorkerHandler)
	r.Get("/etas", api.IndexHandler)

	// vapidpr, vapidpu, _ := vapidkeys.Generate()
	// fmt.Println(vapidpr)
	// fmt.Println(vapidpu)

	vapidPrivateKey = api.cfg.PrivateVapidKey

	r.Post("/sendNotification", api.SendNotificationHandler)

	// iTRAK data feed endpoint
	r.Get("/datafeed", api.DataFeedHandler)

	api.handler = r

	return &api, nil
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		ListenURL:    "0.0.0.0:8080",
		Authenticate: true,
	}
	v.SetDefault("api.listenurl", cfg.ListenURL)
	v.SetDefault("api.casurl", cfg.CasURL)
	v.SetDefault("api.authenticate", cfg.Authenticate)
	return cfg
}

func (api *API) Run() {
	if err := http.ListenAndServe(api.cfg.ListenURL, api.handler); err != nil {
		log.WithError(err).Error("Unable to serve.")
	}
}

func (api *API) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()
	keys := make([]string, 0, len(r.Form))
	for k := range r.Form {
		keys = append(keys, k)
	}
	data := Request{}
	json.Unmarshal([]byte(keys[0]), &data)
	sub := Subscript{}
	json.Unmarshal([]byte(data.Subscription), &sub)
	d := time.Duration(data.Delay) * time.Millisecond
	s := &webpush.Subscription{}
	json.Unmarshal([]byte(sub.Subscription), s)
	time.Sleep(d)
	_, err := webpush.SendNotification([]byte(data.Campus), s, &webpush.Options{
		Subscriber:			"shuttletrackertest@gmail.com",
		VAPIDPublicKey: 	vapidPublicKey,
		VAPIDPrivateKey: 	vapidPrivateKey,
		TTL:				300,
	})
	if (err != nil){
		// Error handling
	}
}

func (api *API) ServiceWorkerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/src/serviceworker.js")
}

// IndexHandler serves the index page.
func (api *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// AdminHandler serves the admin page.
func (api *API) AdminHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query()
	//authenticated with a ticket, redirect the request
	if len(id["ticket"]) != 0 {
		http.Redirect(w, r, "/admin", 301)
	}
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, "static/admin.html")
}

//KeyHandler sends Mapbox api key to authenticated user
func (api *API) KeyHandler(w http.ResponseWriter, r *http.Request) {
	err := WriteJSON(w, api.cfg.MapboxAPIKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
